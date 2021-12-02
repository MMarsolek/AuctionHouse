import { createContext, useContext, useState } from 'react';
import { useCookies } from 'react-cookie';
import biddrClient from './biddrClient/biddrClient'

const AuthContext = createContext(null);

const AuthProvider = ({children}) => {
    const [user, setUser] = useState();
    const [cookies, setCookie, removeCookie] = useCookies(['token', 'username']);
    
    const login = async (username, password) => {
        try {
            const response = await biddrClient.userLogIn(username, password);
            setUser(response);
            setCookie('token', response.authToken);
            setCookie('username', response.username);
        } catch (ex) {
            console.log(`Error logging in: ${ex}`)
            return false;
        }

        return true;
    }

    const logout = async () => {
        setUser(null);
        removeCookie('token');
        removeCookie('username');
        await biddrClient.userLogout();
    }

    const refreshSession = async () => {
        if (!cookies.token) {
            return false;
        }

        biddrClient.userAuth = cookies.token;
        const response = await biddrClient.getUserInformation(cookies.username);
        if (response) {
            setUser(response);
            return true;
        }

        setUser(null);
        biddrClient.userAuth = '';
        return false;
    }

    const value = { user, login, logout, refreshSession };
    return <AuthContext.Provider value={ value }>{children}</AuthContext.Provider>;
}

const useAuth = () => useContext(AuthContext);

export { AuthContext, AuthProvider, useAuth };
