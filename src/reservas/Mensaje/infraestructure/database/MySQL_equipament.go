package database

import (
	"fmt"
	"reserva/src/core"

	"log"
)

type MySQLEquipament struct {
	conn *core.Conn_MySQL
}

func NewMySQLEquipament() *MySQLEquipament {
	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQLEquipament{conn: conn}
}

func (mysql *MySQLEquipament) Save(cname string, category string, ccondition string) {
	query := "INSERT INTO equipments (cname, category, ccondition) VALUES (?, ?, ?)"

	result, err := mysql.conn.ExecutePreparedQuery(query, cname, category, ccondition)
	if err != nil {
		log.Fatalf("Error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Equipo guardado: %d", rowsAffected)
	}
}

func (mysql *MySQLEquipament) GetAll() ([]map[string]interface{}, error) {
	query := "SELECT * FROM equipments"
	rows := mysql.conn.FetchRows(query)
	defer rows.Close()

	var equipaments []map[string]interface{}
	for rows.Next() {
		var id int
		var cname, category, ccondition string
		if err := rows.Scan(&id, &cname, &category, &ccondition); err != nil {
			return nil, err
		}
		equipament := map[string]interface{}{
			"id":        id,
			"cname":      cname,
			"category":  category,
			"ccondition": ccondition,
		}
		equipaments = append(equipaments, equipament)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipaments, nil
}

func (mysql *MySQLEquipament) GetById(id int) ([]map[string]interface{}, error) {
	query := "SELECT * FROM equipments WHERE id = ?"
	rows := mysql.conn.FetchRows(query, id)
	if rows == nil {
		return nil, fmt.Errorf("no se pudo ejecutar la consulta o no hay resultados")
	}
	defer rows.Close()

	var equipaments []map[string]interface{}
	for rows.Next() {
		var id int
		var cname, category, ccondition string
		if err := rows.Scan(&id, &cname, &category, &ccondition); err != nil {
			return nil, err
		}
		equipament := map[string]interface{}{
			"id":        id,
			"name":      cname,
			"category":  category,
			"condition": ccondition,
		}
		equipaments = append(equipaments, equipament)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipaments, nil
}

func (mysql *MySQLEquipament) GetCondition(condition string) ([]map[string]interface{}, error){
	query := "SELECT * FROM equipments WHERE ccondition = ?"
	rows := mysql.conn.FetchRows(query, condition)
	defer rows.Close()

	var equipments []map[string]interface{}
	for rows.Next() {
		var id int
		var cname, category, ccondition string
		if err := rows.Scan(&id, &cname, &category, &ccondition); err != nil {
			return nil, err
		}
		equipment := map[string]interface{}{
			"id": id,
			"cname" :  cname,
			"category" : category,
			"ccondition" : ccondition,
		}
		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipments, nil
}


func (mysql *MySQLEquipament) Update(id int, cname string, category string, ccondition string) {
	query := "UPDATE equipments SET cname = ?, category = ?, ccondition = ? WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, cname, category, ccondition, id)
	if err != nil {
		log.Fatalf("Error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Equipo actualizado: %d", rowsAffected)
	}
}

func (mysql *MySQLEquipament) Delete(id int) {
	query := "DELETE FROM equipments WHERE id = ?"
	result, err := mysql.conn.ExecutePreparedQuery(query, id)
	if err != nil {
		log.Fatalf("Error al ejecutar la consulta: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 1 {
		log.Printf("[MySQL] - Equipo eliminado: %d", rowsAffected)
	}
}