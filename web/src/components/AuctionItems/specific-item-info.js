import React, {Component} from 'react'
import './auction-items.css'
import MakeBid from './Bid/make-bid'




export default class SpecificItem extends Component{

    state = {
        clicked: false,
        item: {},
    }

    constructor(props) {
        super(props);
        this.state.item = props.itemInfo;
        this.updateItem = this.updateItem.bind(this);
    }

    handleClick = () => {
        this.setState({ clicked : !this.state.clicked })
    }

    updateItem(itemChanges) {
        this.setState({ item: { ...this.state.item, ...itemChanges }})
    }


    render(){
        return(
            <div className = "all-item-flex">
                <div className="list-items-form-field">
                    <div className="specific-item-list">
                        <ul className= 'unordered-list'>
                            <li className="items-name"  onClick={this.handleClick }> Item Name: {
                                this.state.item.name
                            } 
                            </li>
                            
                            <li className='bid-amount'>
                                Current Bid: {this.state.item.bidAmount}
                            </li>
                            
                            <div className= 'description-and-image'>
                            {
                                this.state.clicked &&
                                    <li className="items-description"> Description: 
                                    {this.state.item.description}
                                    </li>
                                    
                                }
                                {this.state.item.image &&
                                    this.state.clicked  &&
                                <li className="items-image" >
                                    <img src={this.state.item.image} width='300' height='200' alt="" />
                                    <MakeBid itemName={this.state.item.name} updateItem={this.updateItem} />
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