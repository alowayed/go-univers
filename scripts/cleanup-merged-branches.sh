#!/bin/bash

# cleanup-merged-branches.sh
# Automatically detect and clean up local branches that have been squash-merged remotely

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
MAIN_BRANCH="main"
PROTECTED_BRANCHES=("main" "master" "develop" "development")
WORKING_BRANCH_PATTERNS=() # e.g. ("feature/long-running-task" "user/my-wip-branch")

echo -e "${BLUE}🧹 Branch Cleanup Tool${NC}"
echo "Scanning for branches that have been squash-merged..."
echo

# Ensure we're on main and up to date
echo -e "${BLUE}📡 Updating main branch...${NC}"
git checkout "$MAIN_BRANCH" > /dev/null 2>&1
git fetch --all --prune > /dev/null 2>&1
git pull origin "$MAIN_BRANCH" > /dev/null 2>&1

mapfile -t local_branches < <(git branch --format='%(refname:short)' | grep -v "^$MAIN_BRANCH$")

# Arrays to track branches
merged_branches=()
working_branches=()
unknown_branches=()

# Associative arrays to cache git operation results
declare -A branch_commits_not_in_main
declare -A branch_remote_exists

echo -e "${BLUE}🔍 Analyzing branches...${NC}"

for branch in "${local_branches[@]}"; do
    # Skip if it's a protected branch
    if [[ " ${PROTECTED_BRANCHES[@]} " =~ " ${branch} " ]]; then
        continue
    fi
    
    # Skip if it's a known working branch
    is_working=false
    for pattern in "${WORKING_BRANCH_PATTERNS[@]}"; do
        if [[ "$branch" == $pattern ]]; then
            is_working=true
            working_branches+=("$branch")
            break
        fi
    done
    
    if [[ "$is_working" == true ]]; then
        continue
    fi
    
    # Check if remote tracking branch exists (cache result)
    if git ls-remote --exit-code --heads origin "$branch" > /dev/null 2>&1; then
        branch_remote_exists["$branch"]="true"
        remote_exists="true"
    else
        branch_remote_exists["$branch"]="false"
        remote_exists="false"
    fi
    
    # Check if all commits from this branch are in main (cache result)
    commits_not_in_main=$(git cherry "$MAIN_BRANCH" "$branch" 2>/dev/null | grep "^+" | wc -l)
    branch_commits_not_in_main["$branch"]="$commits_not_in_main"
    
    # Check if there are any differences in file content
    if git diff --quiet "$MAIN_BRANCH"..."$branch" --; then
        file_diff_count=0
    else
        file_diff_count=1
    fi
    
    # Enhanced detection for ecosystem branches (common pattern)
    is_ecosystem_merged=false
    if [[ "$branch" =~ (feat|feature).*(ecosystem|support) ]]; then
        # Check if the ecosystem directory exists in main
        # Filter for actual directories (4th field and beyond, not files in pkg/ecosystem/)
        ecosystem_dirs=$(git diff --name-only "$MAIN_BRANCH"..."$branch" 2>/dev/null | grep "pkg/ecosystem/" | awk -F'/' 'NF >= 4 {print $3}' | sort -u)
        all_ecosystems_exist=true
        
        for ecosystem in $ecosystem_dirs; do
            if [[ -n "$ecosystem" && ! -d "pkg/ecosystem/$ecosystem" ]]; then
                all_ecosystems_exist=false
                break
            fi
        done
        
        # If all ecosystem directories exist in main, likely merged
        if [[ "$all_ecosystems_exist" == true && -n "$ecosystem_dirs" ]]; then
            is_ecosystem_merged=true
        fi
    fi
    
    # Determine branch status with enhanced logic
    if [[ "$remote_exists" == "false" && "$commits_not_in_main" -eq 0 ]] || \
       [[ "$file_diff_count" -eq 0 ]] || \
       [[ "$remote_exists" == "false" && "$is_ecosystem_merged" == true ]]; then
        merged_branches+=("$branch")
    else
        unknown_branches+=("$branch")
    fi
done

# Report findings
echo
echo -e "${GREEN}✅ Working branches (will be kept):${NC}"
if [[ ${#working_branches[@]} -eq 0 ]]; then
    echo "  None found"
else
    for branch in "${working_branches[@]}"; do
        echo "  - $branch"
    done
fi

echo
echo -e "${YELLOW}🔀 Likely merged branches (candidates for deletion):${NC}"
if [[ ${#merged_branches[@]} -eq 0 ]]; then
    echo "  None found"
else
    for branch in "${merged_branches[@]}"; do
        if [[ "${branch_remote_exists["$branch"]}" == "true" ]]; then
            remote_status="(remote exists)"
        else
            remote_status="(remote deleted)"
        fi
        echo "  - $branch $remote_status"
    done
fi

echo
echo -e "${BLUE}❓ Unknown status branches (manual review recommended):${NC}"
if [[ ${#unknown_branches[@]} -eq 0 ]]; then
    echo "  None found"
else
    for branch in "${unknown_branches[@]}"; do
        commits_not_in_main="${branch_commits_not_in_main["$branch"]}"
        if [[ "${branch_remote_exists["$branch"]}" == "true" ]]; then
            remote_status="(remote exists)"
        else
            remote_status="(remote deleted)"
        fi
        echo "  - $branch (${commits_not_in_main} unique commits) $remote_status"
    done
fi

# Offer to delete merged branches
if [[ ${#merged_branches[@]} -gt 0 ]]; then
    echo
    echo -e "${YELLOW}Would you like to delete the likely merged branches? [y/N]${NC}"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo
        echo -e "${GREEN}🗑️  Deleting merged branches...${NC}"
        
        for branch in "${merged_branches[@]}"; do
            echo "  Deleting $branch..."
            git branch -D "$branch" > /dev/null 2>&1
            echo -e "    ${GREEN}✅ Deleted $branch${NC}"
        done
        
        echo
        echo -e "${GREEN}🎉 Cleanup complete! Deleted ${#merged_branches[@]} branches.${NC}"
    else
        echo -e "${BLUE}No branches deleted. Run with individual branch names to delete specific branches.${NC}"
    fi
else
    echo
    echo -e "${GREEN}🎉 No merged branches found. Your repository is clean!${NC}"
fi

# Show final branch status
echo
echo -e "${BLUE}📊 Final branch summary:${NC}"
remaining_branches=$(git branch --format='%(refname:short)' | grep -v "^$MAIN_BRANCH$" | wc -l)
echo "  Total local branches (excluding main): $remaining_branches"

# Usage instructions
echo
echo -e "${BLUE}💡 Usage tips:${NC}"
echo "  - Add working branch patterns to WORKING_BRANCH_PATTERNS in this script"
echo "  - Run this script periodically to keep your local repository clean"
echo "  - Use 'git branch -D <branch>' to manually delete specific branches"