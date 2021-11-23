import React, {Component} from "react";
import { Button } from "../Button";
import {MenuItems} from "./menu-items"
import "./nav-bar.css"
import {button} from "../button"
class Navbar extends Component {
    state = { clicked: false}

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
                <ul className = {this.state.clicked ? "NavMenu active" : "NavMenu"}>
                    {MenuItems.map((item, index) =>{
                        return (
                            <li key = {index}>
                                <a className={item.cName} href={item.url}> 
                                {item.title}
                                </a>
                            </li>
                        )
                    })}
                </ul>
                <Button>Sign up</Button>
            </nav>
        )
    }
}

export default Navbar
