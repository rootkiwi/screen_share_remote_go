#!/bin/bash

sass --style compressed --sourcemap=none --watch web/style.sass:../static/css/style.min.css
