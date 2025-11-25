import React from 'react';
import clsx from 'clsx';

const FeatureBadge = ({ type, children, className }) => {
  const baseClasses = 'feature-badge';

  const typeClasses = {
    enterprise: 'feature-badge--enterprise',
    new: 'feature-badge--new',
    beta: 'feature-badge--beta',
    experimental: 'feature-badge--experimental',
  };

  const icons = {
    enterprise: '🏢',
    new: '✨',
    beta: '🧪',
    experimental: '🔬',
  };

  return (
    <span className={clsx(baseClasses, typeClasses[type], className)}>
      {icons[type]} {children}
    </span>
  );
};

export default FeatureBadge;
