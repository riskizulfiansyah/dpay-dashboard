export const useSidebar = () => {
  const isExpanded = useState('sidebar-expanded', () => true);
  const isOverlay = useState('sidebar-overlay', () => false);

  const checkScreenSize = () => {
    if (typeof window !== 'undefined') {
      isOverlay.value = window.innerWidth < 1024;
    }
  };

  const toggle = () => {
    isExpanded.value = !isExpanded.value;
  };

  const close = () => {
    if (isOverlay.value) {
      isExpanded.value = false;
    }
  };

  return {
    isExpanded,
    isOverlay,
    checkScreenSize,
    toggle,
    close,
  };
};
