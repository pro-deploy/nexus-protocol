import React from 'react';
import clsx from 'clsx';

const ApiMethod = ({ method, children, className }) => {
  const baseClasses = 'api-method';

  const methodClasses = {
    GET: 'api-method--get',
    POST: 'api-method--post',
    PUT: 'api-method--put',
    DELETE: 'api-method--delete',
    PATCH: 'api-method--patch',
  };

  return (
    <code className={clsx(baseClasses, methodClasses[method], className)}>
      {method} {children}
    </code>
  );
};

export default ApiMethod;
