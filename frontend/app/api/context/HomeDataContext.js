import { createContext, useContext } from 'react';
const HomeDataContext = createContext(null);
export const useHomeData = () => {
  return useContext(HomeDataContext);
};
export default HomeDataContext;