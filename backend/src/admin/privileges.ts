export const isSuperUser = ({ currentAdmin }) => currentAdmin && currentAdmin.role === 'SUPER_ADMIN';

export const SuperUserAction = {
  list: {
    isAccessible: isSuperUser,
  },
  edit: {
    isAccessible: isSuperUser,
  },
  delete: {
    isAccessible: isSuperUser,
  },
  bulkDelete: {
    isAccessible: isSuperUser,
  },
  new: {
    isAccessible: isSuperUser,
  },
  search: {
    isAccessible: isSuperUser,
  },
  show: {
    isAccessible: isSuperUser,
  },
};
