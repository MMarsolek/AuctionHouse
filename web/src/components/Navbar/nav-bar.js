import { Component } from 'react';
import { Link } from 'react-router-dom';
import { AuthContext } from '../../auth-provider';
import "./nav-bar.css"

class Navbar extends Component {
    static contextType = AuthContext;
    state = { 
        clicked: false,
        loggedIn: false,
    }

    handleClick = () => {
        this.setState({ clicked : !this.state.clicked })
    }

    render(){

        return(
            <nav className = "NavbarItems">
                <h1 className="NavBarLogo">Auction House  <i className="fa fa-store"></i></h1>
                <div className="MenuIcon" onClick={this.handleClick}>
                    <i className = { this.state.clicked ? "fas fa-times" : "fas fa-bars"}></i>  
                </div>
                { this.context.user && <Link to="/auctions">Auctions</Link>}
                { this.context.user?.permission === 'Admin' && <Link to="/createUser">Create User</Link> }
                { this.context.user?.permission === 'Admin' && <Link to="/createItem">Create Item</Link>}
                { this.context.user && <Link to="/" onClick={async () => await this.context.logout()}>Logout</Link>}
            </nav>
        )
    }
}

export default Navbar
