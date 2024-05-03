-- Copyright 2024 The Mothership Authors
-- SPDX-License-Identifier: Apache-2.0

CREATE VIEW batches_view AS
    SELECT
        b.*,
        (SELECT COUNT(*) FROM entries WHERE batch_name = b.name) as entry_count
    FROM
        batches b;