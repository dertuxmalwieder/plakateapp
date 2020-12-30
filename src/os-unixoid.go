// +build linux solaris dragonfly freebsd netbsd openbsd darwin

/*
 * The contents of this file are subject to the terms of the
 * Common Development and Distribution License, Version 1.0 only
 * (the "License").  You may not use this file except in compliance
 * with the License.
 *
 * See the file LICENSE in this distribution for details.
 * A copy of the CDDL is also available via the Internet at
 * http://www.opensource.org/licenses/cddl1.txt
 *
 * When distributing Covered Code, include this CDDL HEADER in each
 * file and include the contents of the LICENSE file from this
 * distribution.
 */

package main

import (
	"log"

	"github.com/sevlyar/go-daemon"
)

func startServer() {
	// Unter Unix, Linux, BSD usw. können wir forken.
	// Machen wir also!
	cntxt := &daemon.Context{
		PidFileName: "plakateapp.pid",
		PidFilePerm: 0644,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Fork fehlgeschlagen: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	fmt.Println("Fork erfolgreich! Die Plakateapp ist jetzt verfügbar.")

	serveHTTP()
}
