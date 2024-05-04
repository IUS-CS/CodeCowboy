#!/bin/bash

CLASS_LIST=$(gh classroom list | tail -n +4)
CLASS_IDS=$(cut -w -f1 <<< "$CLASS_LIST")
CLASS_NAMES=$(cut -w -f2 <<< "$CLASS_LIST")
CLASS_URLS=$(cut -w -f3 <<< "$CLASS_LIST")
