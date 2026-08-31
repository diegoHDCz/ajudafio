-- name: InsertAuditLog :exec
INSERT INTO audit_logs (user_id, action, metadata)
VALUES (@user_id, @action, @metadata);
