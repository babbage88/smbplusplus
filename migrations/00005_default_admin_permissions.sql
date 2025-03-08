-- +goose Up

-- +goose StatementBegin
INSERT INTO public.app_permissions (permission_name, permission_description) VALUES
    ('CreateRole', 'Permission to Create User Roles'),
    ('AlterRole', 'Alter Role properties and Permissions to Create User Roles'),
    ('CreatePermission', 'Create Permission'),
    ('AlterPermission', 'Alter Permission properties')
ON CONFLICT (permission_name) DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO public.app_permissions (permission_name, permission_description) VALUES
    ('DelteUser', 'Delete Users'),
    ('DeleteRole', 'Delete Roles'),
    ('DeltePermission', 'Delete Permissions')
ON CONFLICT (permission_name) DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
DECLARE
    devuser_id uuid;
    dev_adminroleid uuid;
BEGIN
    -- Find the Admin roll uuid for the default Admin user
    SELECT id INTO dev_adminroleid FROM public.user_roles WHERE role_name = 'Admin';
    -- Find the user ID for the username of the default admin user.

    INSERT INTO public.role_permission_mapping (role_id, permission_id)
    SELECT dev_adminroleid, id FROM public.app_permissions
    WHERE permission_name IN ('DeleteUser', 'DeleteRole', 'DeletePermission');

    INSERT INTO public.role_permission_mapping (role_id, permission_id)
    SELECT dev_adminroleid, id FROM public.app_permissions
    WHERE permission_name IN ('CreateRole', 'AlterRole', 'CreatePermission', 'AlterPermission');
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
DECLARE
    dev_adminroleid uuid;
BEGIN
    -- Find the Admin role UUID
    SELECT id INTO dev_adminroleid FROM public.user_roles WHERE role_name = 'Admin';

    -- Remove permission mappings for the added permissions
    DELETE FROM public.role_permission_mapping
    WHERE role_id = dev_adminroleid AND permission_id IN (
        SELECT id FROM public.app_permissions
        WHERE permission_name IN ('DeleteUser', 'DeleteRole', 'DeletePermission',
                                  'CreateRole', 'AlterRole', 'CreatePermission', 'AlterPermission')
    );

    -- Delete the permissions added in the Up migration
    DELETE FROM public.app_permissions
    WHERE permission_name IN ('DeleteUser', 'DeleteRole', 'DeletePermission',
                              'CreateRole', 'AlterRole', 'CreatePermission', 'AlterPermission');
END $$;
-- +goose StatementEnd

