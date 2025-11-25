import React from 'react';
import clsx from 'clsx';

const StatusIndicator = ({ status, children, className }) => {
  const baseClasses = 'status-indicator';

  const statusClasses = {
    success: 'status-indicator--success',
    error: 'status-indicator--error',
    warning: 'status-indicator--warning',
    info: 'status-indicator--info',
    pending: 'status-indicator--pending',
  };

  return (
    <span className={clsx(baseClasses, statusClasses[status], className)}>
      {children}
    </span>
  );
};

export default StatusIndicator;
