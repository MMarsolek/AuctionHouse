import { Component } from 'react';
import { Link, Navigate, useLocation } from 'react-router-dom';
import { AuthContext } from '../../auth-provider';
import "./nav-bar.css"

class Navbar extends Component {
    static contextType = AuthContext;
    state = { 
        clicked: false,
    }

    handleClick = () => {
        this.setState({ clicked : !this.state.clicked })
    }

    async componentDidMount() {
        await this.context.refreshSession();
    }

    handleLogout = async event => {
        event.preventDefault();
        await this.context.logout();
    }

    render(){
        return(
            <nav className = "NavbarItems">
                { !this.context.user && this.props.location.pathname !== '/' && <Navigate to="/"  />}
                <h1 className="NavBarLogo">Auction House  <i className="fa fa-store"></i></h1>
                <div className="MenuIcon" onClick={this.handleClick}>
                    <i className = { this.state.clicked ? "fas fa-times" : "fas fa-bars"}></i>  
                </div>
                { this.context.user && <Link to="/auctions">Auctions</Link>}
                { this.context.user?.permission === 'Admin' && <Link to="/createUser">Create User</Link> }
                { this.context.user?.permission === 'Admin' && <Link to="/createItem">Create Item</Link>}
                { this.context.user && <button onClick={this.handleLogout}>Logout</button>}
            </nav>
        )
    }
}

const withRoutingFields = (Component) => {
    return props => {
        const location = useLocation();
        return <Component {...props} location={location} />
    }
}

export default withRoutingFields(Navbar);
