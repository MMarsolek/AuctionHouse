import { Component } from 'react';
import { Navigate } from 'react-router-dom';
import { AuthContext } from '../../auth-provider';
import biddrClient from '../../biddrClient/biddrClient';
import './auction-items.css';
import SpecificItem from  './specific-item-info';


export default class ListItems extends Component{
    static contextType = AuthContext;

    state = {
        allItems: [],
    }

    async componentDidMount() {
        const itemMap = (await biddrClient.getAllItems()).reduce((prev, cur) => {prev[cur.name] = cur; return prev}, {});
        (await biddrClient.getHighestBidForAll()).forEach(bidInfo => {
            const item = itemMap[bidInfo.item.name];
            item.bidAmount = bidInfo.bidAmount;
            item.bidder = bidInfo.bidder;
        })

        this.setState({allItems: (Object.values(itemMap))})
    };

    render(){
        return(
            <div className = "all-item-flex">
                { !this.context.user && <Navigate to="/" replace={true} />}
                <div className="item-list">
                    <ul className="items-list">{
                        this.state.allItems.map(item => {
                            return(
                                <li key= {item.name}>
                                    <SpecificItem itemInfo={item}/>
                                </li> 
                            )
                        })}
                    </ul>
                </div>
            </div>
        );
    };
}