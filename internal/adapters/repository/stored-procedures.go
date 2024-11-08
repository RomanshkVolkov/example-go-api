package repository

import (
	"fmt"

	"gorm.io/gorm"
)

// stored procedures
const sp_CreateTableToDocument = `CREATE PROCEDURE [dbo].sp_CreateTableToDocument @id BIGINT
AS
BEGIN
    DECLARE @table NVARCHAR(300) = (SELECT [table] FROM documents WHERE id = @id)
    IF (EXISTS (SELECT *
                FROM INFORMATION_SCHEMA.tables
                WHERE TABLE_SCHEMA = 'dbo'
                  AND TABLE_NAME = @table))
        BEGIN
            THROW 50000, 'RECORD_ALREADY_EXIST', 1;
        END
    DECLARE @columns TABLE
                     (
                         name NVARCHAR(300),
                         type NVARCHAR(100)
                     )

    INSERT INTO @columns
    SELECT 'id',
           'BIGINT IDENTITY
                  PRIMARY KEY'

    INSERT INTO @columns
    SELECT field, type_field
    FROM detail_documents
    WHERE document_id = @id

    DECLARE @sql NVARCHAR(MAX) = 'CREATE TABLE ' + @table + ' (';
    SELECT @sql = @sql + name + ' ' + type + ',' FROM @columns;
    SET @sql = LEFT(@sql, LEN(@sql) - 1) + ');';

    EXEC sp_executesql @sql;
END
`
const sp_GetDocumentTableByID = `CREATE PROCEDURE [dbo].sp_GetDocumentTableByID(@id INT)
AS
BEGIN
    DECLARE @table_name NVARCHAR(100)
    DECLARE @document TABLE
                      (
                          id      BIGINT,
                          name    NVARCHAR(300),
                          [table] NVARCHAR(300)
                      )

    INSERT INTO @document (id, name, [table])
    SELECT id, name, [table]
    FROM documents d
    WHERE d.id = @id

    SET @table_name = (SELECT TOP 1 [table] FROM @document)

    DECLARE @data TABLE
                  (
                      uid         NVARCHAR(300),
                      table_name  NVARCHAR(300),
                      table_align NVARCHAR(20)
                  )
    INSERT INTO @data (uid, table_name, table_align) SELECT 'id', 'none', 'none'

    INSERT INTO @data (uid, table_name, table_align)
    SELECT dd.field, dd.document_key, IIF(dd.type_field = 'DECIMAL(10, 2)' OR dd.type_field = 'INT', 'end', 'start')
    FROM detail_documents dd
    WHERE dd.document_id = @id

    DECLARE @sql NVARCHAR(MAX) = 'SELECT '

    SELECT @sql =
           @sql + CONCAT(uid, '''', uid, '''', ',')

    FROM @data
    SET @sql = LEFT(@sql, LEN(@sql) - 1);
    SET @sql = CONCAT(@sql, ' FROM ', @table_name)

    PRINT @sql

    EXEC sp_executesql @sql

    SELECT uid, table_name AS name, table_align AS align FROM @data WHERE uid <> 'id'

END
`
const sp_AppendDocumentData = `CREATE PROCEDURE [dbo].sp_AppendDocumentData @document_id BIGINT, @file NVARCHAR(MAX)
AS
BEGIN
    DECLARE @table_name NVARCHAR(300);
    DECLARE @details TABLE
                     (
                         field        NVARCHAR(300),
                         document_key NVARCHAR(300),
                         type_field   NVARCHAR(300)
                     );
    DECLARE @sql NVARCHAR(MAX);

    SET @table_name = (SELECT [table] FROM documents WHERE id = @document_id);

    INSERT INTO @details (field, document_key, type_field)
    SELECT field, document_key, type_field
    FROM detail_documents
    WHERE document_id = @document_id;

    SET @sql = 'MERGE INTO ' + @table_name + ' AS trg USING (SELECT ';

    SELECT @sql = @sql + CONCAT(field, ',')
    FROM @details

    SET @sql = dbo.fn_DropEndChar(@sql) + ' FROM OPENJSON(@file) WITH (';

    SELECT @sql = @sql + CONCAT(d.field, ' ', d.type_field, ' ''', '$.', REPLACE(d.document_key, ' ', '_'), '''', ',')
    FROM @details d

    SET @sql = dbo.fn_DropEndChar(@sql) + ')) AS src ON ';

    SELECT @sql = @sql + CONCAT('trg.', d.field, ' = src.', d.field, ' AND ') FROM @details d

    SET @sql = dbo.fn_DropEndChars(@sql, 4) + ' WHEN NOT MATCHED THEN INSERT (';

    SELECT @sql = @sql + field + ',' FROM @details;

    SET @sql = dbo.fn_DropEndChar(@sql) + ') VALUES (';

    SELECT @sql = @sql + 'src.' + field + ', ' FROM @details;

    SET @sql = dbo.fn_DropEndChar(@sql) + ');';

    PRINT @sql

    EXEC sp_executesql @sql, N'@file NVARCHAR(MAX)', @file = @file;
END
`

// functions
const FuncDropEndChar = `CREATE FUNCTION fn_DropEndChar(@txt NVARCHAR(MAX))
    RETURNS NVARCHAR(MAX)
AS BEGIN
    RETURN LEFT(@txt, LEN(@txt) - 1)
END`
const FuncDropEndChars = `CREATE FUNCTION fn_DropEndChars(@txt NVARCHAR(MAX), @steps INT)
    RETURNS NVARCHAR(MAX)
AS
BEGIN
    RETURN LEFT(@txt, LEN(@txt) - @steps)
END
`

func ExistSP(db *gorm.DB, nombreSP string) bool {
	var existe int
	err := db.Raw("SELECT COUNT(*) FROM sys.procedures WHERE name = ?", nombreSP).Scan(&existe).Error

	if err != nil {
		fmt.Println("error when verifying the existence of the stored procedure: %w", err)
		return false
	}

	return existe > 0
}

func ExistFunc(db *gorm.DB, nombreFunc string) bool {
	var existe int
	err := db.Raw("SELECT COUNT(*) FROM sys.objects WHERE name = ? AND type = 'FN'", nombreFunc).Scan(&existe).Error

	if err != nil {
		fmt.Println("error when verifying the existence of the function: %w", err)
		return false
	}

	return existe > 0
}

func ExistTable(db *gorm.DB, name string) bool {
	var rowCount int
	err := db.Raw("SELECT COUNT(*) FROM sys.tables WHERE name = ?", name).Scan(&rowCount).Error

	fmt.Println(rowCount)
	if err != nil {
		fmt.Println("error when verifying the existence of the table: %w", err)
		return false
	}

	return rowCount > 0
}

func MigrateProcedures(db *gorm.DB) {
	// db.Exec("DROP PROCEDURE IF EXISTS sp_CreateTableToDocument")
	exist := ExistSP(db, "sp_CreateTableToDocument")
	if !exist {
		db.Exec(sp_CreateTableToDocument)
	}

	exist = ExistSP(db, "sp_GetDocumentTableByID")
	if !exist {
		db.Exec(sp_GetDocumentTableByID)
	}

	exist = ExistSP(db, "sp_AppendDocumentData")
	if !exist {
		db.Exec(sp_AppendDocumentData)
	}

	exist = ExistFunc(db, "fn_DropEndChar")
	if !exist {
		db.Exec(FuncDropEndChar)
	}

	exist = ExistFunc(db, "fn_DropEndChars")
	if !exist {
		db.Exec(FuncDropEndChars)
	}
}
