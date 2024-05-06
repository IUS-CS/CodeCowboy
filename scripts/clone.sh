#!/bin/bash
set -o errexit -o pipefail

CMD_PWD=$(pwd)
CMD="$0"
CMD_DIR="$(cd "$(dirname "$CMD")" && pwd -P)"


# Basic helpers
out() { echo "$(TZ=UTC date +%Y-%m-%dT%H:%M:%SZ): $*"; }
err() { out "$*" 1>&2; }
vrb() { [ ! "$VERBOSE" ] || out "$@"; }
dbg() { [ ! "$DEBUG" ] || err "$@"; }
die() { err "EXIT: $1" && [ "$2" ] && [ "$2" -ge 0 ] && exit "$2" || exit 1; }

# Show help function to be used below
show_help() {
	awk 'NR>1{print} /^(###|$)/{exit}' "$CMD"
	echo "USAGE: $(basename "$CMD") [class] [assignment]"
	echo "ARGS:"
	MSG=$(awk '/^NARGS=-1; while/,/^esac; done/' "$CMD" | sed -e 's/^[[:space:]]*/  /' -e 's/|/, /' -e 's/)//' | grep '^  -')
	EMSG=$(eval "echo \"$MSG\"")
	echo "$EMSG"
}

# Parse command line options (odd formatting to simplify show_help() above)
NARGS=-1; while [ "$#" -ne "$NARGS" ]; do NARGS=$#; case $1 in
	# SWITCHES
	-h|--help)      # This help message
		show_help; exit 1; ;;
	-d|--debug)     # Enable debugging messages (implies verbose)
		DEBUG=$(( DEBUG + 1 )) && VERBOSE="$DEBUG" && shift && echo "#-INFO: DEBUG=$DEBUG (implies VERBOSE=$VERBOSE)"; ;;
	-v|--verbose)   # Enable verbose messages
		VERBOSE=$(( VERBOSE + 1 )) && shift && echo "#-INFO: VERBOSE=$VERBOSE"; ;;
  --dest)         # Destination directory
    shift && DEST="$1" && shift && vrb "#-INFO: DEST=$DEST"; ;;
	*)
		break;
esac; done

[ "$DEST" ] || DEST=pwd
[ $# -gt 0 -a -z "$CLASS" ]  &&  CLASS="$1"  &&  shift
[ $# -gt 0 -a -z "$ASSIGNMENT" ]  &&  ASSIGNMENT="$1"  &&  shift
[ "$CLASS" ] || die "You must provide a class"
[ "$ASSIGNMENT" ] || die "You must provide an assignment"

mkdir -p "$DEST"
cd "$DEST"

gh classroom clone student-repos -a $(gh classroom assignments -c $(gh classroom list | tail -n +4 | grep "$CLASS" | cut -w -f1)|tail -n +4 | grep "$ASSIGNMENT" | cut -w -f1)

cd "$CMD_PWD"