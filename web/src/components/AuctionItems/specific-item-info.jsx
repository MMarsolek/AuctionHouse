import { Component } from 'react';
import Modal from 'react-modal';
import MakeBid from './Bid/make-bid';
import biddrClient from '../../biddrClient/biddrClient';
import { AuthContext } from '../../auth-provider';
import { withNotifications } from '../../utils';
import './auction-items.scss';

const UpdateItemModal = withNotifications(class UpdateItemModal extends Component {
    state = {
        description: '',
        image: '',
    }

    constructor(props) {
        super(props);

        this.state.description = props.item.description;
        this.state.image = props.item.image;
    }

    handleDescriptionChange = event => {
        this.setState({description: event.target.value});
    };
    handleImageChange = event => {
        this.setState({image: event.target.value});
    };

    handleUpdateClick = async event => {
        event.preventDefault();
        try {
            await biddrClient.updateItem(this.state.description, this.state.image, this.props.item.name);
            this.props.setItemModalState(false);
            this.props.updateItem({description: this.state.description, image: this.state.image});
        } catch (ex) {
            this.props.notify(`Updating item failed: ${ex}`, 'error');
        }
    }

    handleDeleteClick = async event => {
        event.preventDefault();
        try {
            await biddrClient.deleteItem(this.props.item.name);
            this.props.setItemModalState(false);
            this.props.notify('Item deleted, refresh to view changes', 'info');
        } catch (ex) {
            this.props.notify(`Deleting item failed: ${ex}`, 'error');
        }
    }

    render() {
        return (
            <Modal isOpen={this.props.isOpen} className='modal-open'>

                <div className='modal-page'>         
                    <div className="image form-field">
                        <label className="image">Image Link</label>
                        <input  type="text" onChange={this.handleImageChange} className="form-input" placeholder="Image Link"  value={this.state.image} />
                    </div>
                    <div className= "description form-field">
                    <label className="description">Description</label>
                        <textarea  type="text" name="name" onChange={this.handleDescriptionChange} rows= '2' col='300' className= "form-input"
                        placeholder="Description" value={this.state.description} />
                    </div>


                    <button onClick={this.handleUpdateClick} className='button'>Update Item</button>
                    <button onClick={this.handleDeleteClick}className='button'>Delete Item</button>
                    <button onClick={() => this.props.setItemModalState(false)}className='button'>Close</button>
                </div>
            </Modal>
        );
    }
});

class SpecificItem extends Component {
    static contextType = AuthContext;
    state = {
        clicked: false,
        item: {},
        itemModalOpen: false,
    }

    constructor(props) {
        super(props);
        this.state.item = props.itemInfo;
        this.updateItem = this.updateItem.bind(this);
        this.setItemModalState = this.setItemModalState.bind(this);
    }

    handleClick = () => {
        this.setState({ clicked : !this.state.clicked })
    }

    updateItem(itemChanges) {
        this.setState({ item: { ...this.state.item, ...itemChanges }})
    }

    setItemModalState(isOpen) {
        this.setState({itemModalOpen: isOpen})
    }

    render(){
        return(
            <div className="all-item-flex">

                <div className="list-items-form-field">
                    <div className="specific-item-list">
                        <ul className= 'unordered-list'>
                            <li className= {this.state.clicked ? 'name-is-clicked' : 'items-name' }  onClick={this.handleClick }>{this.state.item.name}</li>
                            
                            <li className='bid-amount'>
                                Current Bid Amount: {this.state.item.bidAmount || 0}
                                <MakeBid itemName={this.state.item.name} updateItem={this.updateItem} />

                            </li>
                            
                            <div className= 'description-and-image'>
                            { this.state.clicked &&
                                <li className='button-container'>{ this.context.user.permission === 'Admin' && <button onClick={() => this.setItemModalState(true)} className='edit-item-button'>Edit Item</button> }
                                <UpdateItemModal 
                                    isOpen={this.state.itemModalOpen}
                                    setItemModalState={this.setItemModalState}
                                    item={this.state.item}
                                    updateItem={this.updateItem}
                                /> 
                                </li>
                                }
                            {
                                this.state.clicked &&
                                    <li className="items-description">
                                        {this.state.item.description || 'No description'}
                                    </li>
                            }
                            
                                {this.state.item.image &&
                                    this.state.clicked  &&
                                <li className="items-image" >
                                    <img src={this.state.item.image} width='300' height='200' alt="" />
                                </li>    
                            }
                            
                            </div>

                        </ul>
                    </div>
                </div>
            </div>
        );
    };
}

export default withNotifications(SpecificItem);