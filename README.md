# hsLink
[![Build Status](https://travis-ci.org/YuShuanHsieh/hugo-langs-symlink.svg?branch=master)](https://travis-ci.org/YuShuanHsieh/hugo-langs-symlink)
[![codecov](https://codecov.io/gh/YuShuanHsieh/hugo-langs-symlink/branch/master/graph/badge.svg)](https://codecov.io/gh/YuShuanHsieh/hugo-langs-symlink)

A tool to create symlinks for hugo multi-langs sites.

## Description
hsink tool helps you to create symbolic links for these untranslated pages of Hugo multi-lang site, so that user can see the page with the default language instead of page missing.

## Usage

### Subcommand
1. `create`(c): create symlinks
2. `remove`(r): remove symlinks

### Params

|Param| Description| Example |
|--|--|--|
|dir|the path of content folder|`--dir=$PWD/content`|
|langs|supported langs|`--langs="zh" --langs="zh-tw"`|
|skips|the folders you want to skip|`--skips="iamges"`|
|ext|target file extension|`--ext=.md` (default value is `.md`)|

## Example
```shell
--dir=/home/user/site/content --langs="zh" --langs="zh-tw" --skips="images" create
```
