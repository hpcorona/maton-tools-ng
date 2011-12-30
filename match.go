package main

import "strings"

type MatchPattern []string

func NewMatchPattern(pattern string) MatchPattern {
  return tokenizePathAsArray(pattern)
}

func (mp MatchPattern) Match(name string, caseSensitive bool) bool {
  if mp == nil || len(mp) == 0 {
    return true
  }

  nameArr := tokenizePathAsArray(name)
  return matchPath(mp, nameArr, caseSensitive)
}

func match(pattern, str string, caseSensitive bool) bool {
  patArr := []uint8(pattern)
  strArr := []uint8(str)
  patIdxStart := 0
  patIdxEnd := len(pattern) - 1
  strIdxStart := 0
  strIdxEnd := len(str) - 1
  var ch byte

  containsStar :=  strings.Contains(pattern, "*")
  if !containsStar {
    // No '*'s, so we make a shortcut
    if patIdxEnd != strIdxEnd {
      return false // Pattern and string do not have the same size
    }
    for i := 0; i <= patIdxEnd; i++ {
      ch = patArr[i]
      if ch != '?' {
        if different(caseSensitive, ch, strArr[i]) {
          return false // Character mismatch
        }
      }
    }
    return true // String matches against pattern
  }

  if patIdxEnd == 0 {
    return true // Pattern contains only '*', wich means anything
  }

  // Process characters before first star
  for {
    ch = patArr[patIdxStart]
    if ch == '*' || strIdxStart > strIdxEnd {
      break
    }
    if ch != '?' {
      if different(caseSensitive, ch, strArr[strIdxStart]) {
        return false // Character mismatch
      }
    }
    patIdxStart++
    strIdxStart++
  }

  if strIdxStart > strIdxEnd {
    // All characters in the string are used. Check if only '*'s are
    // left in the pattern. If so, we succeeded. Otherwise failure.
    return allStars(patArr, patIdxStart, patIdxEnd)
  }

  // Process characters after the last star
  for {
    ch = patArr[patIdxEnd]
    if ch == '*' || strIdxStart > strIdxEnd {
      break
    }
    if ch != '?' {
      if different(caseSensitive, ch, strArr[strIdxEnd]) {
        return false
      }
    }
    patIdxEnd--
    strIdxEnd--
  }

  if strIdxStart > strIdxEnd {
    // All characters in the string are used. Check if only '*'s are
    // left in the pattern. If so, we succeeded. Otherwise failure.
    return allStars(patArr, patIdxStart, patIdxEnd)
  }

  // process pattern between stars. padIdxStart and patIdxEnd point
  // always to a '*'.
  for patIdxStart != patIdxEnd && strIdxStart <= strIdxEnd {
    patIdxTmp := -1
    for i := patIdxStart + 1; i <= patIdxEnd; i++ {
      if patArr[i] == '*' {
        patIdxTmp = i
        break
      }
    }
    if patIdxTmp == patIdxStart + 1 {
      // Two stars next to each other, skip the first one.
      patIdxStart++
      continue
    }
    // Find the pattern between padIdxStart & padIdxTmp in str between
    // strIdxStart & strIdxEnd
    patLength := patIdxTmp - patIdxStart - 1
    strLength := strIdxEnd - strIdxStart + 1
    foundIdx := -1

    strLoop := false
    for i := 0; i < strLength - patLength; i++ {
      strLoop = false
      for j := 0; j < patLength; j++ {
        ch = patArr[patIdxStart + j + 1]
        if ch != '?' {
          if different(caseSensitive, ch,
              strArr[strIdxStart + i + j]) {
            strLoop = true
            break
          }
        }
      }
      if strLoop == true {
        continue
      }

      foundIdx = strIdxStart + i
      break
    }

    if foundIdx == -1 {
      return false
    }

    patIdxStart = patIdxTmp
    strIdxStart = foundIdx + patLength
  }

  // All characters in the string are used. Check if only '*'s are left
  // in the pattern. If so, we succeeded. Otherwise failure.
  return allStars(patArr, patIdxStart, patIdxEnd)
}

func allStars(chars []uint8, start, end int) bool {
  for i := start; i <= end; i++ {
    if chars[i] != '*' {
      return false
    }
  }
  return true
}

func different(caseSensitive bool, ch, other uint8) bool {
  if caseSensitive {
    return ch != other
  }

  str1 := string([]uint8{ch})
  str2 := string([]uint8{other})


  str1 = strings.ToUpper(str1)
  str2 = strings.ToUpper(str2)

  return str1 != str2
}

func matchPath(tokenizedPattern, strDirs []string, isCaseSensitive bool) bool {
  patIdxStart := 0
  patIdxEnd := len(tokenizedPattern) - 1
  strIdxStart := 0
  strIdxEnd := len(strDirs) - 1

  // up first '**'
  for patIdxStart <= patIdxEnd && strIdxStart <= strIdxEnd {
    patDir := tokenizedPattern[patIdxStart]
    if patDir == "**" {
      break
    }
    if !match(patDir, strDirs[strIdxStart], isCaseSensitive) {
      return false
    }
    patIdxStart++
    strIdxStart++
  }
  if strIdxStart > strIdxEnd {
    for i := patIdxStart; i <= patIdxEnd; i++ {
      if tokenizedPattern[i] != "**" {
        return false
      }
    }
    return true
  } else {
    if patIdxStart > patIdxEnd {
      // String not exhausted, but pattern is. Failure.
      return false
    }
  }

  // up to last '**'
  for patIdxStart <= patIdxEnd && strIdxStart <= strIdxEnd {
    patDir := tokenizedPattern[patIdxEnd]
    if patDir == "**" {
      break
    }
    if !match(patDir, strDirs[strIdxEnd], isCaseSensitive) {
      return false
    }
    patIdxEnd--
    strIdxEnd--
  }
  if strIdxStart > strIdxEnd {
    // String is exhausted
    for i := patIdxStart; i <= patIdxEnd; i++ {
      if tokenizedPattern[i] != "**" {
        return false
      }
    }
    return true
  }

  for patIdxStart != patIdxEnd && strIdxStart <= strIdxEnd {
    patIdxTmp := -1
    for i := patIdxStart + 1; i <= patIdxEnd; i++ {
      if tokenizedPattern[i] == "**" {
        patIdxTmp = i
        break
      }
    }
    if patIdxTmp == patIdxStart + 1 {
      // '**/**' situation, so skip one
      patIdxStart++
      continue
    }
    // Find the pattern between patIdxStart & patIdxTmp in str between
    // strIdxStart & strIdxEnd
    patLength := patIdxTmp - patIdxStart - 1
    strLength := strIdxEnd - strIdxStart + 1
    foundIdx := -1
    strLoop := false
    for i := 0; i <= strLength - patLength; i++ {
      strLoop = false
      for j := 0; j < patLength; j++ {
        subPat := tokenizedPattern[patIdxStart + j + 1]
        subStr := strDirs[strIdxStart + i + j]
        if !match(subPat, subStr, isCaseSensitive) {
          strLoop = true
          break
        }
      }
      if strLoop {
        continue
      }

      foundIdx = strIdxStart + i
      break
    }

    if foundIdx == -1 {
      return false
    }

    patIdxStart = patIdxTmp
    strIdxStart = foundIdx + patLength
  }

  for i := patIdxStart; i <= patIdxEnd; i++ {
    if tokenizedPattern[i] != "**" {
      return false
    }
  }

  return true
}

func tokenizePathAsArray(path string) []string {
  root := ""
  pathArr := []uint8(path)
  /*
  if fs_isAbs(path) {
    s := dissect(path)
    root = s[0]
    path = s[1]
  }
  */

  var sep uint8 = '/'
  start := 0
  length := len(path)
  count := 0
  for pos := 0; pos < length; pos++ {
    if pathArr[pos] == sep {
      if pos != start {
        count++
      }
      start = pos + 1
    }
  }
  if length != start {
    count++
  }
  var l []string
  if root == "" {
    l = make([]string, count)
  } else {
    l = make([]string, count + 1)
  }

  if root != "" {
    l[0] = root
    count = 1
  } else {
    count = 0
  }
  start = 0
  for pos := 0; pos < length; pos++ {
    if pathArr[pos] == sep {
      if pos != start {
        l[count] = path[start : pos]
        count++
      }
      start = pos + 1
    }
  }
  if length != start {
    l[count] = path[start : ]
  }
  return l
}

