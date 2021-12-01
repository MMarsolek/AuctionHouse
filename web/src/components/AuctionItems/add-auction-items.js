import React, {Component} from 'react'
import biddrClient from '../../biddrClient/biddrClient'
import './auction-items.css'

export default class AddItems extends Component{
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
    };

    render(){
        return(
            <div className = "addition-flex">
                <div className = "add-new-item">
                <form onSubmit={this.handleSubmit} className="add-item-form">
                        <div className= "description form-field">
                        <label className="description"><span className="hidden">Description</span></label>
                            <input type="text" name="name" onChange={this.handleDescriptionChange} className= "form-input"
                            placeholder="Description" />
                        </div>

                        <div className="image form-field">
                            <label className="image"><span className="hidden">Image Link</span></label>
                            <input id="image-link" type="text"  onChange={this.handleImageChange} className="form-input" placeholder="Image Link" />
                        </div>

                        <div className="name form-field">
                            <label className="item-name"><span className="hidden">Name</span></label>
                            <input type="text" name="name" onChange={this.handleNameChange} className= "form-input" placeholder="Item Name" required/>
                        </div>
                        <div className="form-field">
                            <input type="submit" value="Add Item"/>
                        </div>
                    </form>
                </div>
            </div>
        );
    };
}
