INSERT INTO roles(id, name, description)
VALUES
(1, 'admin', 'platform manager'),
(2, 'mentor', 'course manager'),
(3, 'student', 'registered users');

INSERT INTO permissions(id, name, description)
VALUES
(1, 'create', 'resources creation permission'),
(2, 'readOwn', 'own resources read/get permission'),
(3, 'readAll', 'all resources read/get permission'),
(4, 'updateOwn', 'own resources update permission'),
(5, 'updateAll', 'all resources update permission'),
(6, 'deleteOwn', 'own resources deletion permission'),
(7, 'deleteAll', 'own resources deletion permission');

INSERT INTO resources(id, name, description)
VALUES
(1, 'student', 'registered student data'),
(2, 'mentor', 'course manager data'),
(3, 'course_category', 'course category data'),
(4, 'course', 'course data'),
(5, 'course_request', 'course request data');

INSERT INTO rbac(role_id, permission_id, resource_id)
VALUES
(1, 3, 1), -- admin can read all student's data
(1, 5, 2), -- admin can update all mentor data
(1, 7, 2), -- admin can delete all mentor data
(1, 3, 2), -- admin can read all mentor data
(1, 1, 3), -- admin can create course category data
(1, 5, 3), -- admin can update all course category data
(2, 4, 2), -- mentor can update their own data
(3, 4, 1), -- student can update their own data
(2, 1, 4), -- mentor can create course data
(2, 6, 4), -- mentor can delete his own course
(2, 2, 4), -- mentor can read their own course
(2, 4, 4), -- mentor can update their own course
(1, 1, 2), -- admin can create mentor data
(2, 5, 5), -- mentor can update course request data
(2, 2, 5); -- mentor can read their own course-related request