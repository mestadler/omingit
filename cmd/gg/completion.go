package main

import (
	"fmt"
	"os"
)

func completionCmd(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: gg completion <bash|zsh|fish>\n")
		return 1
	}
	switch args[0] {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "fish":
		fmt.Print(fishCompletion)
	default:
		fmt.Fprintf(os.Stderr, "gg completion: unknown shell %q\nSupported: bash, zsh, fish\n", args[0])
		return 1
	}
	return 0
}

const bashCompletion = `# gg completion for bash
# Source this file: source <(gg completion bash)

_gg_platform_subcmds() {
	local cmd="$1" cur="$2"
	local subs
	case "$cmd" in
		pr)       subs="create list view status checkout merge close reopen ready review checks comment diff edit" ;;
		issue)    subs="create list view status close reopen comment edit develop transfer lock unlock delete" ;;
		repo)     subs="clone create fork view list rename sync delete edit archive unarchive" ;;
		release)  subs="create list view delete edit upload download" ;;
		label)    subs="create list clone edit delete" ;;
		workflow) subs="list run enable disable view" ;;
		run)      subs="list view watch cancel rerun download" ;;
		gist)     subs="create list view edit delete clone rename" ;;
		project)  subs="list view create close delete edit item-add item-create item-edit item-list mark-template" ;;
		variable) subs="list set delete get" ;;
		secret)   subs="list set delete" ;;
		milestone) subs="list view create edit delete" ;;
		login)    subs="" ;;
	esac
	COMPREPLY=($(compgen -W "$subs" -- "$cur"))
}

_gg_delegate_flags() {
	local cmd="$1" sub="$2" cur="$3"
	local cli="" help_out flags

	command -v gh &>/dev/null && cli="gh"
	[[ -z "$cli" ]] && command -v tea &>/dev/null && cli="tea"

	if [[ -n "$cli" ]]; then
		help_out=$("$cli" "$cmd" "$sub" --help 2>/dev/null)
		if [[ -n "$help_out" ]]; then
			flags=$(echo "$help_out" | sed -n 's/.*\(--[a-z][a-z0-9-]\{1,\}\).*/\1/p' | sort -u)
			if [[ -n "$flags" ]]; then
				COMPREPLY=($(compgen -W "$flags" -- "$cur"))
				return
			fi
		fi
	fi

	COMPREPLY=($(compgen -f -- "$cur"))
}

_gg() {
	local cur="${COMP_WORDS[COMP_CWORD]}"
	local prev="${COMP_WORDS[COMP_CWORD-1]}"

	local platform_cmds="pr issue repo release label milestone login workflow run gist project variable secret"
	local git_cmds="commit push pull status log diff add branch checkout merge rebase stash tag fetch remote clone init config reset restore switch mv rm bisect grep show cherry-pick revert blame clean gc archive bundle notes submodule worktree"
	local all_cmds="$platform_cmds $git_cmds"

	if [[ "$prev" == "-C" ]]; then
		COMPREPLY=($(compgen -d -S/ -- "$cur"))
		return
	fi

	local plat_cmd="" plat_sub=""
	local i
	for (( i=1; i < COMP_CWORD; i++ )); do
		if [[ "${COMP_WORDS[i]}" == "-C" ]]; then
			(( i++ ))
			continue
		fi
		if [[ -n "$plat_cmd" ]]; then
			plat_sub="${COMP_WORDS[i]}"
			break
		fi
		if [[ " $platform_cmds " =~ " ${COMP_WORDS[i]} " ]]; then
			plat_cmd="${COMP_WORDS[i]}"
		fi
	done

	if [[ -n "$plat_sub" ]]; then
		_gg_delegate_flags "$plat_cmd" "$plat_sub" "$cur"
		return
	fi

	if [[ -n "$plat_cmd" ]]; then
		_gg_platform_subcmds "$plat_cmd" "$cur"
		return
	fi

	if [[ " $all_cmds " =~ " $prev " ]]; then
		COMPREPLY=($(compgen -f -- "$cur"))
		return
	fi

	COMPREPLY=($(compgen -W "--version -v --help -h completion -C $platform_cmds $git_cmds" -- "$cur"))
}

complete -F _gg gg
`

const zshCompletion = `#compdef gg

# gg completion for zsh
# Source this in your .zshrc or copy to a fpath directory.

_gg_platform_subcmds() {
	local cmd="$1" cur="$2"
	local subs
	case "$cmd" in
		pr)       subs="create list view status checkout merge close reopen ready review checks comment diff edit" ;;
		issue)    subs="create list view status close reopen comment edit develop transfer lock unlock delete" ;;
		repo)     subs="clone create fork view list rename sync delete edit archive unarchive" ;;
		release)  subs="create list view delete edit upload download" ;;
		label)    subs="create list clone edit delete" ;;
		workflow) subs="list run enable disable view" ;;
		run)      subs="list view watch cancel rerun download" ;;
		gist)     subs="create list view edit delete clone rename" ;;
		project)  subs="list view create close delete edit item-add item-create item-edit item-list mark-template" ;;
		variable) subs="list set delete get" ;;
		secret)   subs="list set delete" ;;
		milestone) subs="list view create edit delete" ;;
		login)    subs="" ;;
	esac
	compadd -- ${=subs}
}

_gg_delegate_flags() {
	local cmd="$1" sub="$2" cur="$3"
	local cli="" help_out flags

	command -v gh &>/dev/null && cli="gh"
	[[ -z "$cli" ]] && command -v tea &>/dev/null && cli="tea"

	if [[ -n "$cli" ]]; then
		help_out=$("$cli" "$cmd" "$sub" --help 2>/dev/null)
		if [[ -n "$help_out" ]]; then
			flags=$(echo "$help_out" | sed -n 's/.*\(--[a-z][a-z0-9-]\{1,\}\).*/\1/p' | sort -u)
			if [[ -n "$flags" ]]; then
				compadd -- ${=flags}
				return
			fi
		fi
	fi

	_files
}

_gg() {
	local cur="${words[CURRENT]}"
	local prev="${words[CURRENT-1]}"

	local platform_cmds="pr issue repo release label milestone login workflow run gist project variable secret"
	local git_cmds="commit push pull status log diff add branch checkout merge rebase stash tag fetch remote clone init config reset restore switch mv rm bisect grep show cherry-pick revert blame clean gc archive bundle notes submodule worktree"
	local all_cmds="$platform_cmds $git_cmds"

	if [[ "$prev" == "-C" ]]; then
		_files -/
		return
	fi

	local plat_cmd="" plat_sub=""
	local i=1
	while (( i < CURRENT )); do
		if [[ "${words[i]}" == "-C" ]]; then
			(( i += 2 ))
			continue
		fi
		if [[ -n "$plat_cmd" ]]; then
			plat_sub="${words[i]}"
			break
		fi
		if [[ " $platform_cmds " =~ " ${words[i]} " ]]; then
			plat_cmd="${words[i]}"
		fi
		(( i++ ))
	done

	if [[ -n "$plat_sub" ]]; then
		_gg_delegate_flags "$plat_cmd" "$plat_sub" "$cur"
		return
	fi

	if [[ -n "$plat_cmd" ]]; then
		_gg_platform_subcmds "$plat_cmd" "$cur"
		return
	fi

	if [[ " $all_cmds " =~ " $prev " ]]; then
		_files
		return
	fi

	compadd -- ${=platform_cmds} ${=git_cmds} --version -v --help -h completion -C
}

_gg "$@"
`

const fishCompletion = `# gg completion for fish
# Source this file: source (gg completion fish | psub)
# Or install: gg completion fish > ~/.config/fish/completions/gg.fish

set -l platform_cmds pr issue repo release label milestone login workflow run gist project variable secret
set -l git_cmds commit push pull status log diff add branch checkout merge rebase stash tag fetch remote clone init config reset restore switch mv rm bisect grep show cherry-pick revert blame clean gc archive bundle notes submodule worktree

# First-argument commands
for cmd in $platform_cmds
    complete -c gg -n "not __fish_seen_subcommand_from $platform_cmds $git_cmds" -a $cmd -d "Platform: $cmd"
end
for cmd in $git_cmds
    complete -c gg -n "not __fish_seen_subcommand_from $platform_cmds $git_cmds" -a $cmd -d "git $cmd"
end

# Special flags and subcommands
complete -c gg -n "not __fish_seen_subcommand_from $platform_cmds $git_cmds" -a "--version" -d "Print version"
complete -c gg -n "not __fish_seen_subcommand_from $platform_cmds $git_cmds" -a "-v" -d "Print version"
complete -c gg -n "not __fish_seen_subcommand_from $platform_cmds $git_cmds" -a "--help" -d "Show help"
complete -c gg -n "not __fish_seen_subcommand_from $platform_cmds $git_cmds" -a "-h" -d "Show help"
complete -c gg -n "not __fish_seen_subcommand_from $platform_cmds $git_cmds" -a "completion" -d "Generate shell completion"

# -C flag: directory argument
complete -c gg -s C -d "Run as if gg was started in <path>" -x -a "(__fish_complete_directories)" \
    -n "not __fish_seen_subcommand_from $platform_cmds $git_cmds"

# pr subcommands
complete -c gg -n "__fish_seen_subcommand_from pr" -a create -d "Create a pull request"
complete -c gg -n "__fish_seen_subcommand_from pr" -a list -d "List pull requests"
complete -c gg -n "__fish_seen_subcommand_from pr" -a view -d "View a pull request"
complete -c gg -n "__fish_seen_subcommand_from pr" -a status -d "Show status of relevant pull requests"
complete -c gg -n "__fish_seen_subcommand_from pr" -a checkout -d "Check out a pull request"
complete -c gg -n "__fish_seen_subcommand_from pr" -a merge -d "Merge a pull request"
complete -c gg -n "__fish_seen_subcommand_from pr" -a close -d "Close a pull request"
complete -c gg -n "__fish_seen_subcommand_from pr" -a reopen -d "Reopen a pull request"
complete -c gg -n "__fish_seen_subcommand_from pr" -a ready -d "Mark a pull request as ready for review"
complete -c gg -n "__fish_seen_subcommand_from pr" -a review -d "Review a pull request"
complete -c gg -n "__fish_seen_subcommand_from pr" -a checks -d "Show checks status"
complete -c gg -n "__fish_seen_subcommand_from pr" -a comment -d "Create a new comment"
complete -c gg -n "__fish_seen_subcommand_from pr" -a diff -d "View changes in a pull request"
complete -c gg -n "__fish_seen_subcommand_from pr" -a edit -d "Edit a pull request"

# issue subcommands
complete -c gg -n "__fish_seen_subcommand_from issue" -a create -d "Create an issue"
complete -c gg -n "__fish_seen_subcommand_from issue" -a list -d "List issues"
complete -c gg -n "__fish_seen_subcommand_from issue" -a view -d "View an issue"
complete -c gg -n "__fish_seen_subcommand_from issue" -a status -d "Show status of relevant issues"
complete -c gg -n "__fish_seen_subcommand_from issue" -a close -d "Close an issue"
complete -c gg -n "__fish_seen_subcommand_from issue" -a reopen -d "Reopen an issue"
complete -c gg -n "__fish_seen_subcommand_from issue" -a comment -d "Comment on an issue"
complete -c gg -n "__fish_seen_subcommand_from issue" -a edit -d "Edit an issue"
complete -c gg -n "__fish_seen_subcommand_from issue" -a develop -d "Manage linked branches"
complete -c gg -n "__fish_seen_subcommand_from issue" -a transfer -d "Transfer issue to another repository"
complete -c gg -n "__fish_seen_subcommand_from issue" -a lock -d "Lock issue conversation"
complete -c gg -n "__fish_seen_subcommand_from issue" -a unlock -d "Unlock issue conversation"
complete -c gg -n "__fish_seen_subcommand_from issue" -a delete -d "Delete an issue"

# repo subcommands
complete -c gg -n "__fish_seen_subcommand_from repo" -a clone -d "Clone a repository"
complete -c gg -n "__fish_seen_subcommand_from repo" -a create -d "Create a new repository"
complete -c gg -n "__fish_seen_subcommand_from repo" -a fork -d "Fork a repository"
complete -c gg -n "__fish_seen_subcommand_from repo" -a view -d "View a repository"
complete -c gg -n "__fish_seen_subcommand_from repo" -a list -d "List repositories"
complete -c gg -n "__fish_seen_subcommand_from repo" -a rename -d "Rename a repository"
complete -c gg -n "__fish_seen_subcommand_from repo" -a sync -d "Sync a repository"
complete -c gg -n "__fish_seen_subcommand_from repo" -a delete -d "Delete a repository"
complete -c gg -n "__fish_seen_subcommand_from repo" -a edit -d "Edit repository settings"
complete -c gg -n "__fish_seen_subcommand_from repo" -a archive -d "Archive a repository"
complete -c gg -n "__fish_seen_subcommand_from repo" -a unarchive -d "Unarchive a repository"

# release subcommands
complete -c gg -n "__fish_seen_subcommand_from release" -a create -d "Create a release"
complete -c gg -n "__fish_seen_subcommand_from release" -a list -d "List releases"
complete -c gg -n "__fish_seen_subcommand_from release" -a view -d "View a release"
complete -c gg -n "__fish_seen_subcommand_from release" -a delete -d "Delete a release"
complete -c gg -n "__fish_seen_subcommand_from release" -a edit -d "Edit a release"
complete -c gg -n "__fish_seen_subcommand_from release" -a upload -d "Upload assets to a release"
complete -c gg -n "__fish_seen_subcommand_from release" -a download -d "Download release assets"

# label subcommands
complete -c gg -n "__fish_seen_subcommand_from label" -a create -d "Create a label"
complete -c gg -n "__fish_seen_subcommand_from label" -a list -d "List labels"
complete -c gg -n "__fish_seen_subcommand_from label" -a clone -d "Clone labels from another repo"
complete -c gg -n "__fish_seen_subcommand_from label" -a edit -d "Edit a label"
complete -c gg -n "__fish_seen_subcommand_from label" -a delete -d "Delete a label"

# workflow subcommands
complete -c gg -n "__fish_seen_subcommand_from workflow" -a list -d "List workflows"
complete -c gg -n "__fish_seen_subcommand_from workflow" -a run -d "Run a workflow"
complete -c gg -n "__fish_seen_subcommand_from workflow" -a enable -d "Enable a workflow"
complete -c gg -n "__fish_seen_subcommand_from workflow" -a disable -d "Disable a workflow"
complete -c gg -n "__fish_seen_subcommand_from workflow" -a view -d "View a workflow"

# run subcommands
complete -c gg -n "__fish_seen_subcommand_from run" -a list -d "List workflow runs"
complete -c gg -n "__fish_seen_subcommand_from run" -a view -d "View a workflow run"
complete -c gg -n "__fish_seen_subcommand_from run" -a watch -d "Watch a workflow run"
complete -c gg -n "__fish_seen_subcommand_from run" -a cancel -d "Cancel a workflow run"
complete -c gg -n "__fish_seen_subcommand_from run" -a rerun -d "Rerun a workflow run"
complete -c gg -n "__fish_seen_subcommand_from run" -a download -d "Download workflow run artifacts"

# gist subcommands
complete -c gg -n "__fish_seen_subcommand_from gist" -a create -d "Create a gist"
complete -c gg -n "__fish_seen_subcommand_from gist" -a list -d "List gists"
complete -c gg -n "__fish_seen_subcommand_from gist" -a view -d "View a gist"
complete -c gg -n "__fish_seen_subcommand_from gist" -a edit -d "Edit a gist"
complete -c gg -n "__fish_seen_subcommand_from gist" -a delete -d "Delete a gist"
complete -c gg -n "__fish_seen_subcommand_from gist" -a clone -d "Clone a gist"
complete -c gg -n "__fish_seen_subcommand_from gist" -a rename -d "Rename a gist"

# project subcommands
complete -c gg -n "__fish_seen_subcommand_from project" -a list -d "List projects"
complete -c gg -n "__fish_seen_subcommand_from project" -a view -d "View a project"
complete -c gg -n "__fish_seen_subcommand_from project" -a create -d "Create a project"
complete -c gg -n "__fish_seen_subcommand_from project" -a close -d "Close a project"
complete -c gg -n "__fish_seen_subcommand_from project" -a delete -d "Delete a project"
complete -c gg -n "__fish_seen_subcommand_from project" -a edit -d "Edit a project"
complete -c gg -n "__fish_seen_subcommand_from project" -a item-add -d "Add an item to a project"
complete -c gg -n "__fish_seen_subcommand_from project" -a item-create -d "Create a draft item"
complete -c gg -n "__fish_seen_subcommand_from project" -a item-edit -d "Edit a project item"
complete -c gg -n "__fish_seen_subcommand_from project" -a item-list -d "List project items"
complete -c gg -n "__fish_seen_subcommand_from project" -a mark-template -d "Mark a project as a template"

# variable subcommands
complete -c gg -n "__fish_seen_subcommand_from variable" -a list -d "List variables"
complete -c gg -n "__fish_seen_subcommand_from variable" -a set -d "Set a variable"
complete -c gg -n "__fish_seen_subcommand_from variable" -a delete -d "Delete a variable"
complete -c gg -n "__fish_seen_subcommand_from variable" -a get -d "Get a variable"

# secret subcommands
complete -c gg -n "__fish_seen_subcommand_from secret" -a list -d "List secrets"
complete -c gg -n "__fish_seen_subcommand_from secret" -a set -d "Set a secret"
complete -c gg -n "__fish_seen_subcommand_from secret" -a delete -d "Delete a secret"

# milestone subcommands
complete -c gg -n "__fish_seen_subcommand_from milestone" -a list -d "List milestones"
complete -c gg -n "__fish_seen_subcommand_from milestone" -a view -d "View a milestone"
complete -c gg -n "__fish_seen_subcommand_from milestone" -a create -d "Create a milestone"
complete -c gg -n "__fish_seen_subcommand_from milestone" -a edit -d "Edit a milestone"
complete -c gg -n "__fish_seen_subcommand_from milestone" -a delete -d "Delete a milestone"

# login
complete -c gg -n "__fish_seen_subcommand_from login" -a "" -d "Authenticate with a platform"
`
