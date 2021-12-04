import React, {Component} from 'react'
import biddrClient from '../../biddrClient/biddrClient'
import { withNotifications } from '../../utils'
import './auction-items.scss'
import '../UserSignIn/user-login'


export default withNotifications(class AddItems extends Component{
    state = {
        description: '',
        image: '',
        name: '',
    };

    handleDescriptionChange = event => {
        this.setState({description: event.target.value});
    };
    handleImageChange = event => {
        this.setState({image: event.target.value});
    };
    handleNameChange = event => {
        this.setState({name: event.target.value});
    };
    
    
    handleSubmit = event => {
        event.preventDefault();
        biddrClient.createNewItem(this.state.description, this.state.image, this.state.name)
        this.props.notify('Item created', 'info');
    };

    render(){
        return(
            <div className = "login-flex">
                <div className="logo"> Create Item</div>
                <div className="login-container">
                <form onClick={this.handleSubmit} className="login-form">
                        <div className="name form-field">
                        <label className="label">Description</label>
                            <input type="text" name="name" onChange={this.handleDescriptionChange} className="form-input"
                            placeholder="Description" />
                        </div>

                        <div className="name form-field">
                            <label className="label">Image Link</label>
                            <input type="text"  onChange={this.handleImageChange} className="form-input" placeholder="Image Link" />
                        </div>

                        <div className="name form-field">
                            <label className="label">Name</label>
                            <input type="text" name="name" onChange={this.handleNameChange} className="form-input" placeholder="Item Name" required/>
                        </div>
                        <div className="button-box">
                            <input type="button" value="Add Item"className="create-button"/>
                        </div>
                    </form>
                </div>
            </div>
        );
    };
});
