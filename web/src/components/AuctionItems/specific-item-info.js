import React, {Component} from 'react'
import './auction-items.css'
import MakeBid from './Bid/make-bid'




export default class SpecificItem extends Component{

    state = {clicked: false}

    handleClick = () => {
        this.setState({ clicked : !this.state.clicked })
    }
    

    render(){
        return(
            <div className = "all-item-flex">
                <div className="list-items-form-field">
                    <div className="specific-item-list">
                        <ul className= 'unordered-list'>    
                            <li className="items-name"  onClick={this.handleClick }> Item Name: {
                                this.props.itemInfo['name'] 
                            } 
                            </li>
                            
                            <li className='bid-amount'>
                                Current Bid: {this.props.itemInfo['bidAmount']}
                            </li>
                            
                            <div className= 'description-and-image'>
                            {
                                this.state.clicked &&
                                    <li className="items-description"> Description: 
                                    {this.props.itemInfo['description']}
                                    </li>
                                    
                                }
                                {this.props.itemInfo['image'] &&
                                    this.state.clicked  &&
                                <li className="items-image" >
                                    <img src= {this.props.itemInfo['image']} width='300' height='200' />
                                    <MakeBid itemName={this.props.itemInfo['name']}/>
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