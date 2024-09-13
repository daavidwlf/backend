package main

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"time"
)

func parseLogs(logs io.Reader) ([]string, error) {
	var logsArray []string

	scanner := bufio.NewScanner(logs)

	for scanner.Scan() {
		logsArray = append(logsArray, scanner.Text())
	}

	err := scanner.Err()

	if err != nil {
		return nil, errors.New("uable to scan logs of io Stream: " + err.Error())
	}

	return logsArray, nil
}

const (
	maxAttempts      = 10
	maxAttemptsPerIP = 10              // max attempts allowed
	blockDuration    = 1 * time.Minute // block duration after too many attempts
	windowDuration   = 1 * time.Minute // time window to reset attempt count
	cleanupInterval  = 1 * time.Minute // interval for cleanup of expired entries
)

var loginAttempts = make(map[string]*loginAttemptInfo)
var mu sync.Mutex

var lastCleanup = time.Now()

func cleanupOldAttempts(now time.Time) {
	for id, info := range loginAttempts {
		if now.Sub(info.lastAttempt) > windowDuration {
			delete(loginAttempts, id)
		}
	}
	lastCleanup = now
}

func trackLoginAttempt(username, ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	if now.Sub(lastCleanup) > cleanupInterval {
		cleanupOldAttempts(now)
	}

	info, exists := loginAttempts[username]
	if !exists {
		info = &loginAttemptInfo{
			attemptCount: 0,
			lastAttempt:  now,
			ipAttempts:   make(map[string]int),
		}
		loginAttempts[username] = info
	}

	if now.Sub(info.lastAttempt) > windowDuration {
		info.attemptCount = 0
	}

	if info.blockedUntil.After(now) {
		return true
	}

	info.attemptCount++
	info.lastAttempt = now

	info.ipAttempts[ip]++
	if info.attemptCount >= maxAttempts {
		info.blockedUntil = now.Add(blockDuration)
	}

	if info.ipAttempts[ip] >= maxAttemptsPerIP {
		return true
	}

	return info.blockedUntil.After(now)
}
