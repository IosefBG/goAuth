package services

//this is a nth after everything works ig but to declare the structure ig
//
//import (
//	"database/sql"
//	"fmt"
//	"log"
//)
//
//package main
//
//import (
//"database/sql"
//"fmt"
//"log"
//)
//
//// CheckPermission checks if a user has the specified permission
//func CheckPermission(db *sql.DB, userID int, permission string) (bool, error) {
//	query := `
//    SELECT 1
//    FROM user_roles ur
//    JOIN role_permissions rp ON ur.role_id = rp.role_id
//    JOIN permissions p ON rp.permission_id = p.id
//    WHERE ur.user_id = $1 AND p.name = $2;
//	`
//	var exists int
//	err := db.QueryRow(query, userID, permission).Scan(&exists)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			// Permission not found
//			return false, nil
//		}
//		return false, err
//	}
//	return exists == 1, nil
//}
//
//func main() {
//	// Assume db is your database connection
//	var db *sql.DB // Replace with actual db connection
//
//	// Simulate user check (e.g., guest user with ID 4)
//	userID := 4
//	permission := "PLACE_ORDER"
//
//	hasPermission, err := CheckPermission(db, userID, permission)
//	if err != nil {
//		log.Fatalf("Error checking permission: %v", err)
//	}
//
//	if hasPermission {
//		fmt.Println("User has permission to place an order.")
//	} else {
//		fmt.Println("User does NOT have permission to place an order.")
//	}
//}
