INSERT INTO users (
    created_at,
    updated_at,
    email,
    profile_image_url,
    volleynet_user_id,
    volleynet_user,
    role,
    pw_salt,
    pw_hash,
    pw_iterations
)
VALUES (
    :created_at,
    :updated_at,
    :email,
    :profile_image_url,
    :volleynet_user_id,
    :volleynet_user,
    :role,
    :pw_salt,
    :pw_hash,
    :pw_iterations
)
RETURNING id