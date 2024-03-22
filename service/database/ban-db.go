package database

import "errors"

func (db *appdbimpl) BanUser(userId uint, targetUserId uint) error {
	_, err := db.c.Exec("INSERT INTO Bans (bannedUserId, userId) VALUES (?, ?)", targetUserId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) UnbanUser(userId uint, targetUserId uint) error {
	res, err := db.c.Exec("DELETE FROM Bans WHERE bannedUserId = ? AND userId = ?", targetUserId, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("target user does not exist in the ban list")
	}

	return nil
}
