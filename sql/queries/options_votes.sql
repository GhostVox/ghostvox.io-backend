-- name: GetVotesByOptionAndPollID :many
SELECT
    votes.*,
    op.*,
    po.*
FROM
    votes
    JOIN options op ON votes.option_id = op.id
    JOIN polls po ON votes.poll_id = po.id
WHERE
    votes.option_id = $1
    AND votes.poll_id = $2;

-- name: CountVotesByOptionAndPollID :many
WITH vote_counts AS (
    SELECT
        option_id,
        COUNT(*) as vote_count
    FROM
        votes
    WHERE
        votes.poll_id = $1
    GROUP BY
        votes.option_id
)
SELECT
    option.name,
    COALESCE(vc.vote_count, 0) as vote_count
FROM
    options option
    LEFT JOIN vote_counts vc ON option.id = vc.option_id
WHERE
    option.poll_id = $2;
