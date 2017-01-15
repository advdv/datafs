// Copyright 2016 Keybase Inc. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// These tests all do multiple operations while a user is unstaged.

package test

import (
	"testing"
	"time"
)

// bob and alice both write(to the same file),
func TestCrConflictWriteFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "world"),
			read(crname("a/b", bob), "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "world"),
			read(crname("a/b", bob), "uh oh"),
		),
	)
}

// bob and alice both create the same entry with different types
func TestCrConflictCreateWithDifferentTypes(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			mkdir("a/c"),
		),
		as(bob, noSync(),
			mkfile("a/c", ""),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "DIR",
				crnameEsc("c", bob): "FILE"}),
			read("a/b", "hello"),
			lsdir("a/c", m{}),
			read(crname("a/c", bob), ""),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "DIR",
				crnameEsc("c", bob): "FILE"}),
			read("a/b", "hello"),
			lsdir("a/c", m{}),
			read(crname("a/c", bob), ""),
		),
	)
}

// bob and alice both create the same file with different types
func TestCrConflictCreateFileWithDifferentTypes(t *testing.T) {
	test(t,
		skip("dokan", "Does not work with Dokan."),
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			mkfile("a/c", ""),
		),
		as(bob, noSync(),
			link("a/c", "b"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "FILE",
				crnameEsc("c", bob): "SYM"}),
			read("a/b", "hello"),
			read("a/c", ""),
			read(crname("a/c", bob), "hello"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "FILE",
				crnameEsc("c", bob): "SYM"}),
			read("a/b", "hello"),
			read("a/c", ""),
			read(crname("a/c", bob), "hello"),
		),
	)
}

// bob and alice both create the same symlink with different contents
func TestCrConflictCreateSymlinkWithDifferentContents(t *testing.T) {
	test(t,
		skip("dokan", "Does not work with Dokan."),
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkfile("a/b", "hello"),
			mkfile("a/c", "world"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			link("a/d", "b"),
		),
		as(bob, noSync(),
			link("a/d", "c"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "FILE", "d$": "SYM",
				crnameEsc("d", bob): "SYM"}),
			read("a/b", "hello"),
			read("a/c", "world"),
			read("a/d", "hello"),
			read(crname("a/d", bob), "world"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "FILE", "d$": "SYM",
				crnameEsc("d", bob): "SYM"}),
			read("a/b", "hello"),
			read("a/c", "world"),
			read("a/d", "hello"),
			read(crname("a/d", bob), "world"),
		),
	)
}

// bob and alice both write(to the same file), but on a non-default day.
func TestCrConflictWriteFileWithAddTime(t *testing.T) {
	timeInc := 25 * time.Hour
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			addTime(timeInc),
			write("a/b", "world"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE",
				crnameAtTimeEsc("b", bob, timeInc): "FILE"}),
			read("a/b", "world"),
			read(crnameAtTime("a/b", bob, timeInc), "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE",
				crnameAtTimeEsc("b", bob, timeInc): "FILE"}),
			read("a/b", "world"),
			read(crnameAtTime("a/b", bob, timeInc), "uh oh"),
		),
	)
}

// bob and alice both write(to the same file),
func TestCrConflictWriteFileWithExtension(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/foo.tar.gz", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/foo.tar.gz", "world"),
		),
		as(bob, noSync(),
			write("a/foo.tar.gz", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"foo.tar.gz$": "FILE", crnameEsc("foo.tar.gz", bob): "FILE"}),
			read("a/foo.tar.gz", "world"),
			read(crname("a/foo.tar.gz", bob), "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"foo.tar.gz$": "FILE", crnameEsc("foo.tar.gz", bob): "FILE"}),
			read("a/foo.tar.gz", "world"),
			read(crname("a/foo.tar.gz", bob), "uh oh"),
		),
	)
}

// bob and alice both create the same file
func TestCrConflictCreateFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "world"),
			read(crname("a/b", bob), "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "world"),
			read(crname("a/b", bob), "uh oh"),
		),
	)
}

// alice setattr's a file, while bob removes, recreates and writes to
// a file of the same name. Regression test for KBFS-668.
func TestCrConflictSetattrVsRecreatedFileInRoot(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			setex("a", true),
		),
		as(bob, noSync(),
			write("a", "uh oh"),
			rm("a"),
			mkfile("a", "world"),
			reenableUpdates(),
			lsdir("", m{"a$": "EXEC", crnameEsc("a", bob): "FILE"}),
			read("a", "hello"),
			read(crname("a", bob), "world"),
		),
		as(alice,
			lsdir("", m{"a$": "EXEC", crnameEsc("a", bob): "FILE"}),
			read("a", "hello"),
			read(crname("a", bob), "world"),
		),
	)
}

// bob creates a directory with the same name that alice used for a file
func TestCrConflictCauseRenameOfMergedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			write("a/b/c", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", alice): "FILE"}),
			read(crname("a/b", alice), "world"),
			read("a/b/c", "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", alice): "FILE"}),
			read(crname("a/b", alice), "world"),
			read("a/b/c", "uh oh"),
		),
	)
}

// bob creates a directory with the same name that alice used for a
// file that used to exist at that location
func TestCrConflictCauseRenameOfMergedRecreatedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			rm("a/b"),
			write("a/b/c", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", alice): "FILE"}),
			read(crname("a/b", alice), "world"),
			read("a/b/c", "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", alice): "FILE"}),
			read(crname("a/b", alice), "world"),
			read("a/b/c", "uh oh"),
		),
	)
}

// bob renames a file over one modified by alice.
func TestCrConflictUnmergedRenameFileOverModifiedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
			write("a/c", "world"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "uh oh"),
		),
		as(bob, noSync(),
			rename("a/c", "a/b"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "uh oh"),
			read(crname("a/b", bob), "world"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "uh oh"),
			read(crname("a/b", bob), "world"),
		),
	)
}

// bob modifies and renames a file that was modified by alice.
func TestCrConflictUnmergedRenameModifiedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			rename("a/b", "a/c"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "FILE"}),
			read("a/b", "world"),
			read("a/c", "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "FILE"}),
			read("a/b", "world"),
			read("a/c", "uh oh"),
		),
	)
}

// bob modifies and renames a file that was modified by alice, while
// alice also made a file with the new name.
func TestCrConflictUnmergedRenameModifiedFileAndConflictFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
			mkfile("a/c", "CONFLICT"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			rename("a/b", "a/c"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "FILE", crnameEsc("c", bob): "FILE"}),
			read("a/b", "world"),
			read("a/c", "CONFLICT"),
			read(crname("a/c", bob), "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "FILE", crnameEsc("c", bob): "FILE"}),
			read("a/b", "world"),
			read("a/c", "CONFLICT"),
			read(crname("a/c", bob), "uh oh"),
		),
	)
}

// bob modifies and renames (to another dir) a file that was modified
// by alice.
func TestCrConflictUnmergedRenameAcrossDirsModifiedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			rename("a/b", "b/c"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE"}),
			read("a/b", "world"),
			lsdir("b/", m{"c$": "FILE"}),
			read("b/c", "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE"}),
			read("a/b", "world"),
			lsdir("b/", m{"c$": "FILE"}),
			read("b/c", "uh oh"),
		),
	)
}

// bob sets the mtime on and renames a file that had its mtime set by alice.
func TestCrConflictUnmergedRenameSetMtimeFile(t *testing.T) {
	targetMtime1 := time.Now().Add(1 * time.Minute)
	targetMtime2 := targetMtime1.Add(1 * time.Minute)
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			setmtime("a/b", targetMtime1),
		),
		as(bob, noSync(),
			setmtime("a/b", targetMtime2),
			rename("a/b", "a/c"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "FILE"}),
			mtime("a/b", targetMtime1),
			mtime("a/c", targetMtime2),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "FILE"}),
			mtime("a/b", targetMtime1),
			mtime("a/c", targetMtime2),
		),
	)
}

// bob renames a file from a new directory over one modified by alice.
func TestCrConflictUnmergedRenameFileInNewDirOverModifiedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
			write("a/c", "world"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "uh oh"),
		),
		as(bob, noSync(),
			rename("a/c", "e/c"),
			rename("e/c", "a/b"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			lsdir("e/", m{}),
			read("a/b", "uh oh"),
			read(crname("a/b", bob), "world"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			lsdir("e/", m{}),
			read("a/b", "uh oh"),
			read(crname("a/b", bob), "world"),
		),
	)
}

// bob renames an existing directory over one created by alice.
// TODO: it would be better if this weren't a conflict.
func TestCrConflictUnmergedRenamedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b/c", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/d/e", "world"),
		),
		as(bob, noSync(),
			write("a/b/f", "uh oh"),
			rename("a/b", "a/d"),
			reenableUpdates(),
			lsdir("a/", m{"d$": "DIR", crnameEsc("d", bob): "DIR"}),
			lsdir("a/d", m{"e": "FILE"}),
			lsdir(crname("a/d", bob), m{"c": "FILE", "f": "FILE"}),
			read(crname("a/d", bob)+"/c", "hello"),
			read("a/d/e", "world"),
			read(crname("a/d", bob)+"/f", "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"d$": "DIR", crnameEsc("d", bob): "DIR"}),
			lsdir("a/d", m{"e": "FILE"}),
			lsdir(crname("a/d", bob), m{"c": "FILE", "f": "FILE"}),
			read(crname("a/d", bob)+"/c", "hello"),
			read("a/d/e", "world"),
			read(crname("a/d", bob)+"/f", "uh oh"),
		),
	)
}

// bob renames a directory over one made non-empty by alice
func TestCrConflictUnmergedRenameDirOverNonemptyDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a/b"),
			mkfile("a/c/d", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			mkfile("a/b/e", "uh oh"),
		),
		as(bob, noSync(),
			rename("a/c", "a/b"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", bob): "DIR"}),
			lsdir("a/b", m{"e": "FILE"}),
			lsdir(crname("a/b", bob), m{"d": "FILE"}),
		),
		as(alice,
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", bob): "DIR"}),
			lsdir("a/b", m{"e": "FILE"}),
			lsdir(crname("a/b", bob), m{"d": "FILE"}),
		),
	)
}

// alice renames an existing directory over one created by bob. TODO:
// it would be better if this weren't a conflict.
func TestCrConflictMergedRenamedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b/c", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b/f", "uh oh"),
			rename("a/b", "a/d"),
		),
		as(bob, noSync(),
			write("a/d/e", "world"),
			reenableUpdates(),
			lsdir("a/", m{"d$": "DIR", crnameEsc("d", bob): "DIR"}),
			lsdir("a/d", m{"c": "FILE", "f": "FILE"}),
			read("a/d/c", "hello"),
			read(crname("a/d", bob)+"/e", "world"),
			read("a/d/f", "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"d$": "DIR", crnameEsc("d", bob): "DIR"}),
			lsdir("a/d", m{"c": "FILE", "f": "FILE"}),
			read("a/d/c", "hello"),
			read(crname("a/d", bob)+"/e", "world"),
			read("a/d/f", "uh oh"),
		),
	)
}

// alice renames a file over one modified by bob.
func TestCrConflictMergedRenameFileOverModifiedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
			write("a/c", "world"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/c", "a/b"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "world"),
			read(crname("a/b", bob), "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "world"),
			read(crname("a/b", bob), "uh oh"),
		),
	)
}

// alice modifies and renames a file that was modified by bob.
func TestCrConflictMergedRenameModifiedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
			rename("a/b", "a/c"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "FILE"}),
			read("a/b", "uh oh"),
			read("a/c", "world"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "FILE"}),
			read("a/b", "uh oh"),
			read("a/c", "world"),
		),
	)
}

// alice modifies and renames a file that was modified by bob, while
// bob also made a file with the new name.
func TestCrConflictMergedRenameModifiedFileAndConflictFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "uh oh"),
			rename("a/b", "a/c"),
		),
		as(bob, noSync(),
			write("a/b", "world"),
			mkfile("a/c", "CONFLICT"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "FILE", crnameEsc("c", bob): "FILE"}),
			read("a/b", "world"),
			read("a/c", "uh oh"),
			read(crname("a/c", bob), "CONFLICT"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "FILE", crnameEsc("c", bob): "FILE"}),
			read("a/b", "world"),
			read("a/c", "uh oh"),
			read(crname("a/c", bob), "CONFLICT"),
		),
	)
}

// alice modifies and renames (to another dir) a file that was modified
// by bob.
func TestCrConflictMergedRenameAcrossDirsModifiedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
			rename("a/b", "b/c"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE"}),
			read("a/b", "uh oh"),
			lsdir("b/", m{"c$": "FILE"}),
			read("b/c", "world"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE"}),
			read("a/b", "uh oh"),
			lsdir("b/", m{"c$": "FILE"}),
			read("b/c", "world"),
		),
	)
}

// alice sets the mtime on and renames a file that had its mtime set by bob.
func TestCrConflictMergedRenameSetMtimeFile(t *testing.T) {
	targetMtime1 := time.Now().Add(1 * time.Minute)
	targetMtime2 := targetMtime1.Add(1 * time.Minute)
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			setmtime("a/b", targetMtime1),
			rename("a/b", "a/c"),
		),
		as(bob, noSync(),
			setmtime("a/b", targetMtime2),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", "c$": "FILE"}),
			mtime("a/b", targetMtime2),
			mtime("a/c", targetMtime1),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", "c$": "FILE"}),
			mtime("a/b", targetMtime2),
			mtime("a/c", targetMtime1),
		),
	)
}

// alice and both both rename(the same file, causing a copy.),
func TestCrConflictRenameSameFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b", "a/c"),
		),
		as(bob, noSync(),
			rename("a/b", "a/d"),
			reenableUpdates(),
			lsdir("a/", m{"c": "FILE", "d": "FILE"}),
			read("a/c", "hello"),
			read("a/d", "hello"),
		),
		as(alice,
			lsdir("a/", m{"c": "FILE", "d": "FILE"}),
			read("a/c", "hello"),
			read("a/d", "hello"),
			write("a/c", "world"),
		),
		as(bob,
			read("a/c", "world"),
			read("a/d", "hello"),
		),
	)
}

// alice and both both rename(the same executable file, causing a copy.),
func TestCrConflictRenameSameEx(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b", "hello"),
			setex("a/b", true),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b", "a/c"),
		),
		as(bob, noSync(),
			rename("a/b", "a/d"),
			reenableUpdates(),
			lsdir("a/", m{"c": "EXEC", "d": "EXEC"}),
			read("a/c", "hello"),
			read("a/d", "hello"),
		),
		as(alice,
			lsdir("a/", m{"c": "EXEC", "d": "EXEC"}),
			read("a/c", "hello"),
			read("a/d", "hello"),
			write("a/c", "world"),
		),
		as(bob,
			read("a/c", "world"),
			read("a/d", "hello"),
		),
	)
}

// alice and both both rename(the same symlink.),
func TestCrConflictRenameSameSymlink(t *testing.T) {
	test(t,
		skip("dokan", "Does not work with Dokan."),
		users("alice", "bob"),
		as(alice,
			write("a/foo", "hello"),
			link("a/b", "foo"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b", "a/c"),
		),
		as(bob, noSync(),
			rename("a/b", "a/d"),
			reenableUpdates(),
			lsdir("a/", m{"foo": "FILE", "c": "SYM", "d": "SYM"}),
			read("a/c", "hello"),
			read("a/d", "hello"),
		),
		as(alice,
			lsdir("a/", m{"foo": "FILE", "c": "SYM", "d": "SYM"}),
			read("a/c", "hello"),
			read("a/d", "hello"),
			write("a/c", "world"),
		),
		as(bob,
			read("a/c", "world"),
			read("a/d", "world"),
		),
	)
}

// alice and bob both rename(the same directory, causing a symlink to),
// be created.
func TestCrConflictRenameSameDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b/c", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b", "a/d"),
		),
		as(bob, noSync(),
			rename("a/b", "a/e"),
			reenableUpdates(),
			lsdir("a/", m{"d": "DIR", "e": "SYM"}),
			read("a/d/c", "hello"),
			read("a/e/c", "hello"),
		),
		as(alice,
			lsdir("a/", m{"d": "DIR", "e": "SYM"}),
			read("a/d/c", "hello"),
			read("a/e/c", "hello"),
			write("a/d/f", "world"),
			read("a/e/f", "world"),
		),
		as(bob,
			read("a/e/f", "world"),
		),
	)
}

// alice and bob both rename(the same directory, causing a symlink to),
// be created.
func TestCrConflictRenameSameDirUpward(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b/c/d/e/foo", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b/c/d/e", "a/e"),
		),
		as(bob, noSync(),
			rename("a/b/c/d/e", "a/b/c/d/f"),
			reenableUpdates(),
			lsdir("a/", m{"b": "DIR", "e": "DIR"}),
			lsdir("a/e", m{"foo": "FILE"}),
			lsdir("a/b/c/d", m{"f": "SYM"}),
			lsdir("a/b/c/d/f", m{"foo": "FILE"}),
			read("a/e/foo", "hello"),
			lsdir("a/b/c/d/f", m{"foo": "FILE"}),
		),
		as(alice,
			lsdir("a/", m{"b": "DIR", "e": "DIR"}),
			lsdir("a/e", m{"foo": "FILE"}),
			lsdir("a/b/c/d", m{"f": "SYM"}),
			lsdir("a/b/c/d/f", m{"foo": "FILE"}),
			read("a/e/foo", "hello"),
			lsdir("a/b/c/d/f", m{"foo": "FILE"}),
			write("a/e/foo2", "world"),
		),
		as(bob,
			read("a/b/c/d/f/foo2", "world"),
		),
	)
}

// alice and bob both rename(the same directory, causing a symlink to),
// be created.
func TestCrConflictRenameSameDirMergedUpward(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b/c/d/e/foo", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b/c/d/e", "a/b/c/d/f"),
		),
		as(bob, noSync(),
			rename("a/b/c/d/e", "a/e"),
			reenableUpdates(),
			lsdir("a/", m{"b": "DIR", "e": "SYM"}),
			lsdir("a/e", m{"foo": "FILE"}),
			lsdir("a/b/c/d", m{"f": "DIR"}),
			lsdir("a/b/c/d/f", m{"foo": "FILE"}),
			read("a/e/foo", "hello"),
			lsdir("a/b/c/d/f", m{"foo": "FILE"}),
		),
		as(alice,
			lsdir("a/", m{"b": "DIR", "e": "SYM"}),
			lsdir("a/e", m{"foo": "FILE"}),
			lsdir("a/b/c/d", m{"f": "DIR"}),
			lsdir("a/b/c/d/f", m{"foo": "FILE"}),
			read("a/e/foo", "hello"),
			lsdir("a/b/c/d/f", m{"foo": "FILE"}),
			write("a/e/foo2", "world"),
		),
		as(bob,
			read("a/b/c/d/f/foo2", "world"),
		),
	)
}

func TestCrConflictRenameSameDirDownward(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b/foo", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b", "a/c/d/e/f"),
		),
		as(bob, noSync(),
			rename("a/b", "a/g"),
			reenableUpdates(),
			lsdir("a/", m{"c": "DIR", "g": "SYM"}),
			lsdir("a/c/d/e/f", m{"foo": "FILE"}),
			lsdir("a/g", m{"foo": "FILE"}),
			read("a/c/d/e/f/foo", "hello"),
			read("a/g/foo", "hello"),
		),
		as(alice,
			lsdir("a/", m{"c": "DIR", "g": "SYM"}),
			lsdir("a/c/d/e/f", m{"foo": "FILE"}),
			lsdir("a/g", m{"foo": "FILE"}),
			read("a/c/d/e/f/foo", "hello"),
			read("a/g/foo", "hello"),
			write("a/c/d/e/f/foo2", "world"),
		),
		as(bob,
			read("a/g/foo2", "world"),
		),
	)
}

func TestCrConflictRenameSameDirSideways(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b/c/d/foo", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b/c/d", "a/e/f/g"),
		),
		as(bob, noSync(),
			rename("a/b/c/d", "a/b/c/h"),
			reenableUpdates(),
			lsdir("a/e/f", m{"g": "DIR"}),
			lsdir("a/b/c", m{"h": "SYM"}),
			lsdir("a/e/f/g", m{"foo": "FILE"}),
			lsdir("a/b/c/h", m{"foo": "FILE"}),
			read("a/e/f/g/foo", "hello"),
			read("a/b/c/h/foo", "hello"),
		),
		as(alice,
			lsdir("a/e/f", m{"g": "DIR"}),
			lsdir("a/b/c", m{"h": "SYM"}),
			lsdir("a/e/f/g", m{"foo": "FILE"}),
			lsdir("a/b/c/h", m{"foo": "FILE"}),
			read("a/e/f/g/foo", "hello"),
			read("a/b/c/h/foo", "hello"),
			write("a/e/f/g/foo2", "world"),
		),
		as(bob,
			read("a/b/c/h/foo2", "world"),
		),
	)
}

// bob renames an existing directory over one created by alice, twice.
// TODO: it would be better if this weren't a conflict.
func TestCrConflictUnmergedRenamedDirDouble(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			write("a/b/c", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/d/e", "world"),
		),
		as(bob, noSync(),
			write("a/b/f", "uh oh"),
			rename("a/b", "a/d"),
			reenableUpdates(),
			lsdir("a/", m{"d$": "DIR", crnameEsc("d", bob): "DIR"}),
			lsdir("a/d", m{"e": "FILE"}),
			lsdir(crname("a/d", bob), m{"c": "FILE", "f": "FILE"}),
			read(crname("a/d", bob)+"/c", "hello"),
			read("a/d/e", "world"),
			read(crname("a/d", bob)+"/f", "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"d$": "DIR", crnameEsc("d", bob): "DIR"}),
			lsdir("a/d", m{"e": "FILE"}),
			lsdir(crname("a/d", bob), m{"c": "FILE", "f": "FILE"}),
			read(crname("a/d", bob)+"/c", "hello"),
			read("a/d/e", "world"),
			read(crname("a/d", bob)+"/f", "uh oh"),
			rm("a/d/e"),
			rm("a/d"),
			write("a/b/c", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/d/e", "world"),
		),
		as(bob, noSync(),
			write("a/b/f", "uh oh"),
			rename("a/b", "a/d"),
			reenableUpdates(),
			lsdir("a/", m{"d$": "DIR", crnameEsc("d", bob) + "$": "DIR", crnameEsc("d", bob) + ` \(1\)`: "DIR"}),
			lsdir("a/d", m{"e": "FILE"}),
			lsdir(crname("a/d", bob)+" (1)", m{"c": "FILE", "f": "FILE"}),
			read("a/d/e", "world"),
		),
		as(alice,
			lsdir("a/", m{"d$": "DIR", crnameEsc("d", bob) + "$": "DIR", crnameEsc("d", bob) + ` \(1\)`: "DIR"}),
			lsdir("a/d", m{"e": "FILE"}),
			lsdir(crname("a/d", bob)+" (1)", m{"c": "FILE", "f": "FILE"}),
			read("a/d/e", "world"),
		),
	)
}

// bob and alice both write(to the same file),
func TestCrConflictWriteFileDouble(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "world"),
			read(crname("a/b", bob), "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "world"),
			read(crname("a/b", bob), "uh oh"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "another write"),
		),
		as(bob, noSync(),
			write("a/b", "uh oh again!"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob) + "$": "FILE", crnameEsc("b", bob) + ` \(1\)`: "FILE"}),
			read("a/b", "another write"),
			read(crname("a/b", bob), "uh oh"),
			read(crname("a/b", bob)+" (1)", "uh oh again!"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob) + "$": "FILE", crnameEsc("b", bob) + ` \(1\)`: "FILE"}),
			read("a/b", "another write"),
			read(crname("a/b", bob), "uh oh"),
			read(crname("a/b", bob)+" (1)", "uh oh again!"),
		),
	)
}

// bob and alice both write(to the same file),
func TestCrConflictWriteFileDoubleWithExtensions(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/file.tar.gz", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/file.tar.gz", "world"),
		),
		as(bob, noSync(),
			write("a/file.tar.gz", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"file.tar.gz$": "FILE", crnameEsc("file.tar.gz", bob): "FILE"}),
			read("a/file.tar.gz", "world"),
			read(crname("a/file.tar.gz", bob), "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"file.tar.gz$": "FILE", crnameEsc("file.tar.gz", bob): "FILE"}),
			read("a/file.tar.gz", "world"),
			read(crname("a/file.tar.gz", bob), "uh oh"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/file.tar.gz", "another write"),
		),
		as(bob, noSync(),
			write("a/file.tar.gz", "uh oh again!"),
			reenableUpdates(),
			lsdir("a/", m{"file.tar.gz$": "FILE", crnameEsc("file.tar.gz", bob) + "$": "FILE", crnameEsc("file", bob) + ` \(1\).tar.gz`: "FILE"}),
			read("a/file.tar.gz", "another write"),
			read(crname("a/file.tar.gz", bob), "uh oh"),
			read(crname("a/file", bob)+" (1).tar.gz", "uh oh again!"),
		),
		as(alice,
			lsdir("a/", m{"file.tar.gz$": "FILE", crnameEsc("file.tar.gz", bob) + "$": "FILE", crnameEsc("file", bob) + ` \(1\).tar.gz`: "FILE"}),
			read("a/file.tar.gz", "another write"),
			read(crname("a/file.tar.gz", bob), "uh oh"),
			read(crname("a/file", bob)+" (1).tar.gz", "uh oh again!"),
		),
	)
}

// bob causes a rename(cycle with a conflict while unstaged),
func TestCrRenameCycleWithConflict(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkdir("a/b"),
			mkdir("a/c"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/c", "a/b/c"),
		),
		as(bob, noSync(),
			rename("a/b", "a/c/b"),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", bob): "FILE"}),
			read(crname("a/b", bob), "uh oh"),
			lsdir("a/b/", m{"c": "DIR"}),
			lsdir("a/b/c", m{"b": "SYM"}),
			lsdir("a/b/c/b", m{"c": "DIR"}),
		),
		as(alice,
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", bob): "FILE"}),
			read(crname("a/b", bob), "uh oh"),
			lsdir("a/b/", m{"c": "DIR"}),
			lsdir("a/b/c", m{"b": "SYM"}),
			lsdir("a/b/c/b", m{"c": "DIR"}),
			write("a/b/d", "hello"),
		),
		as(bob,
			read("a/b/c/b/d", "hello"),
		),
	)
}

// bob causes a rename(cycle with two conflicts while unstaged),
func TestCrRenameCycleWithTwoConflicts(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkdir("a/b"),
			mkdir("a/c"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/c", "a/b/c"),
			write("a/b/c/b", "uh oh"),
		),
		as(bob, noSync(),
			rename("a/b", "a/c/b"),
			write("a/b", "double uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", bob): "FILE"}),
			read(crname("a/b", bob), "double uh oh"),
			lsdir("a/b/", m{"c": "DIR"}),
			lsdir("a/b/c", m{"b$": "SYM", crnameEsc("b", alice): "FILE"}),
			lsdir("a/b/c/b", m{"c": "DIR"}),
		),
		as(alice,
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", bob): "FILE"}),
			read(crname("a/b", bob), "double uh oh"),
			lsdir("a/b/", m{"c": "DIR"}),
			lsdir("a/b/c", m{"b$": "SYM", crnameEsc("b", alice): "FILE"}),
			lsdir("a/b/c/b", m{"c": "DIR"}),
			write("a/b/d", "hello"),
		),
		as(bob,
			read("a/b/c/b/d", "hello"),
		),
	)
}

// bob causes a rename(cycle with two conflicts while unstaged),
func TestCrRenameCycleWithConflictAndMergedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkdir("a/b"),
			mkdir("a/c"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/c", "a/b/c"),
			mkdir("a/b/c/b"),
		),
		as(bob, noSync(),
			rename("a/b", "a/c/b"),
			write("a/b", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", bob): "FILE"}),
			read(crname("a/b", bob), "uh oh"),
			lsdir("a/b/", m{"c": "DIR"}),
			lsdir("a/b/c", m{"b$": "DIR", crnameEsc("b", bob): "SYM"}),
			lsdir(crname("a/b/c/b", bob), m{"c": "DIR"}),
			lsdir("a/b/c/b", m{}),
		),
		as(alice,
			lsdir("a/", m{"b$": "DIR", crnameEsc("b", bob): "FILE"}),
			read(crname("a/b", bob), "uh oh"),
			lsdir("a/b/", m{"c": "DIR"}),
			lsdir("a/b/c", m{"b$": "DIR", crnameEsc("b", bob): "SYM"}),
			lsdir(crname("a/b/c/b", bob), m{"c": "DIR"}),
			lsdir("a/b/c/b", m{}),
			write("a/b/d", "hello"),
		),
		as(bob,
			read(crname("a/b/c/b", bob)+"/d", "hello"),
		),
	)
}

// alice and bob both truncate the same file to different sizes
func TestCrBothTruncateFileDifferentSizes(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			truncate("a/b", 4),
		),
		as(bob, noSync(),
			truncate("a/b", 3),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "hell"),
			read(crname("a/b", bob), "hel"),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", "hell"),
			read(crname("a/b", bob), "hel"),
		),
	)
}

// alice and bob both truncate the same file to different sizes, after
// truncating to the same size
func TestCrBothTruncateFileDifferentSizesAfterSameSize(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			truncate("a/b", 0),
		),
		as(bob, noSync(),
			truncate("a/b", 0),
			truncate("a/b", 3),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", ""),
			read(crname("a/b", bob), string(make([]byte, 3))),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			read("a/b", ""),
			read(crname("a/b", bob), string(make([]byte, 3))),
		),
	)
}

// alice and bob both set the mtime on a file
func TestCrBothSetMtimeFile(t *testing.T) {
	targetMtime1 := time.Now().Add(1 * time.Minute)
	targetMtime2 := targetMtime1.Add(1 * time.Minute)
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			setmtime("a/b", targetMtime1),
		),
		as(bob, noSync(),
			setmtime("a/b", targetMtime2),
			reenableUpdates(),
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			mtime("a/b", targetMtime1),
			mtime(crname("a/b", bob), targetMtime2),
		),
		as(alice,
			lsdir("a/", m{"b$": "FILE", crnameEsc("b", bob): "FILE"}),
			mtime("a/b", targetMtime1),
			mtime(crname("a/b", bob), targetMtime2),
		),
	)
}

// alice and bob both set the mtime on a dir
func TestCrBothSetMtimeDir(t *testing.T) {
	targetMtime1 := time.Now().Add(1 * time.Minute)
	targetMtime2 := targetMtime1.Add(1 * time.Minute)
	test(t,
		skip("dokan", "Dokan can't read mtimes on symlinks."),
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			setmtime("a", targetMtime1),
		),
		as(bob, noSync(),
			setmtime("a", targetMtime2),
			reenableUpdates(),
			lsdir("", m{"a$": "DIR", crnameEsc("a", bob): "SYM"}),
			mtime("a", targetMtime1),
			mtime(crname("a", bob), targetMtime2),
		),
		as(alice,
			lsdir("", m{"a$": "DIR", crnameEsc("a", bob): "SYM"}),
			mtime("a", targetMtime1),
			mtime(crname("a", bob), targetMtime2),
		),
	)
}