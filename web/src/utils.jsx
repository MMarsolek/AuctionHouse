import { useNotifications } from 'reapop';
import { useLocation } from 'react-router';

export const withNotifications = Component => (props) => {
    const {notify} = useNotifications(); 
    return <Component {...props} notify={notify} />
}

export const withRoutingFields = (Component) => props => {
    const location = useLocation();
    return <Component {...props} location={location} />
}
