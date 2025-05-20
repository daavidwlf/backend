package utils

import (
	customTypes "backend/src/types"
	"bufio"
	"errors"
	"io"
	"sync"
	"time"
)

func ParseLogs(logs io.Reader) ([]string, error) {
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

var loginAttempts = make(map[string]*customTypes.LoginAttemptInfo)
var mu sync.Mutex

var lastCleanup = time.Now()

func CleanupOldAttempts(now time.Time) {
	for id, info := range loginAttempts {
		if now.Sub(info.LastAttempt) > windowDuration {
			delete(loginAttempts, id)
		}
	}
	lastCleanup = now
}

func TrackLoginAttempt(username, ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	if now.Sub(lastCleanup) > cleanupInterval {
		CleanupOldAttempts(now)
	}

	info, exists := loginAttempts[username]
	if !exists {
		info = &customTypes.LoginAttemptInfo{
			AttemptCount: 0,
			LastAttempt:  now,
			IpAttempts:   make(map[string]int),
		}
		loginAttempts[username] = info
	}

	if now.Sub(info.LastAttempt) > windowDuration {
		info.AttemptCount = 0
	}

	if info.BlockedUntil.After(now) {
		return true
	}

	info.AttemptCount++
	info.LastAttempt = now

	info.IpAttempts[ip]++
	if info.AttemptCount >= maxAttempts {
		info.BlockedUntil = now.Add(blockDuration)
	}

	if info.IpAttempts[ip] >= maxAttemptsPerIP {
		return true
	}

	return info.BlockedUntil.After(now)
}
