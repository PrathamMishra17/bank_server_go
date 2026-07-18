-- name: TransferAmount :execresult
INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
)VALUES(?,?,?);

-- name: GetTransfer :execresult
SELECT * FROM transfers 
WHERE id = ? LIMIT 1;

-- name: TransferLists :execresult
SELECT * FROM transfers
WHERE from_account_id = ? OR 
to_account_id = ?
ORDER BY id
LIMIT  ?
OFFSET  ?;

