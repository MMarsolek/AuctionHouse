
import React, {Component} from 'react'
import biddrClient from '../../biddrClient/biddrClient'
import './auction-items.css'
import SpecificItem from  './specific-item-info'


export default class ListItems extends Component{
    state = {
        allItems: []
    }
    handleSubmit = async event => {
        event.preventDefault();
            this.setState({allItems: (await biddrClient.getAllItems())})
            // console.log(await biddrClient.getAllItems())
    };


    

    render(){
        return(
            <div className = "all-item-flex">
                <div className = "all-items" onClick={this.handleSubmit} className="get-all-items">
                    <div className="list-items-form-field">
                        <input type="button" value="View All"/>               
                        <div className="item-list">
                            <ul className="items-list">{
                                this.state.allItems.map(item =>
                                   {
                                    return(<li key= {item.name}>
                                        <SpecificItem itemInfo={item}/>
                                        </li> )
                                })}
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        );
    };
}