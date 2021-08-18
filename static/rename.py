#!/usr/bin/env python3
# rename.py version 1.0
# Copyright (C) 2019 Dakota Walsh

import os         # file operations
import sys        # system operations
import argparse   # argument parsing
import hashlib    # calculate hash

def tooBig(old_path, max_size):
	# check is the file is over the max size
	if (os.path.getsize(old_path) <= (max_size*1024*1024)):
		return False
	else:
		return True

def getHash(hold, f_name):
	# return a hash of the file in question with the extension
	root = ''
	ext	 = ''
	# however first we'll get the extension if there is one
	if '.' in f_name:
		root, ext = os.path.splitext(f_name)
	f_path = os.path.join(hold, f_name)
	# calculate the hash of the file
	f_hash = hashlib.sha1(open(f_path,'rb').read()).hexdigest()
	return(f_hash + ext)

def getArgs():
	# get the arguments with argparse
	arg = argparse.ArgumentParser(description="rename.py")
	arg.add_argument('--file', '-f', required=True, help='The file which incrond detected.')
	arg.add_argument('--hold', '-d', required=True, help='The directory which was watched.')
	arg.add_argument('--web', '-w', required=True, help='The web directory.')
	arg.add_argument('--max', '-m', help='The max file size in MB. If no max size is set there will be no max size.')
	return arg.parse_args()

def checkDirs(hold, web):
	# check that the directories exist
	if not (os.path.isdir(hold)):
		print ("ERROR: Hold Directory not found! Check your spelling and ensure the directory exists.")
		sys.exit(1)
	if not (os.path.isdir(web)):
		print ("ERROR: Web Directory not found! Check your spelling and ensure the directory exists.")
		sys.exit(1)

def main():
	# get the passed arguments
	arguments = getArgs()
	old_name  = arguments.file
	hold      = arguments.hold
	web       = arguments.web
	old_path  = os.path.join(hold, old_name)

	# make sure the hold and web directories exist
	checkDirs(hold, web)

	# check max size
	if arguments.max:
		max_size = int(arguments.max)
		if tooBig(old_path, max_size):
			os.remove(old_path)
			sys.exit(0)

	# calculate sum and rename file
	new_name = getHash(hold, old_name) # can't use old_path due to ext check
	new_path = os.path.join(web, new_name)
	os.rename(old_path, new_path)

if __name__ == '__main__':
	main()
