
import ReactDOM from 'react-dom';
import { HashRouter, Routes, Route } from 'react-router-dom';
import NotificationSystem, { NotificationsProvider, atalhoTheme, useNotifications, setUpNotifications } from 'reapop';
import { AuthProvider } from './auth-provider';
import Navbar from './components/Navbar/nav-bar';
import AddItems from './components/AuctionItems/add-auction-items';
import ListItems from './components/AuctionItems/list-auction-items';
import UserLogIn  from './components/UserSignIn/user-login';
import { CookiesProvider } from 'react-cookie';
import UserCreation from './components/UserManagement/user-creation';

const NotificationWrapper = () => {
    const {notifications, dismissNotification} = useNotifications();
    return <NotificationSystem notifications={notifications} dismissNotification={id => dismissNotification(id)} theme={atalhoTheme} />
}

function App() {
    return (
        <CookiesProvider>
            <NotificationsProvider>
                <HashRouter>
                    <AuthProvider>
                        <Navbar />
                        <NotificationWrapper />
                        <Routes>
                            <Route index path="/" element={<UserLogIn />} />
                            <Route path="/auctions" element={<ListItems />} />
                            <Route path="/createUser" element={<UserCreation />} />
                            <Route path="/createItem" element={<AddItems />} />
                        </Routes>
                    </AuthProvider>
                </HashRouter>
            </NotificationsProvider>
        </CookiesProvider>
    )
}

setUpNotifications({
    defaultProps: {
        position: 'top-center',
        dismissible: true,
        dismissAfter: 10000,
    }
});

ReactDOM.render(App(), document.getElementById('app'));
