import axios from 'axios'

class BiddrClient{

    userAuth = '';
    domain = '';

    constructor(baseUrl){
        this.domain = baseUrl
    }
    
    async requestConfig(url, method, {headers, body} = {}){
        const myParameters = {
            'url' : url, 
            'method' : method, 
            'baseURL' : this.domain,  
            'validateStatus' :  status => status < 500
        };
        if(headers){

            myParameters['headers'] = headers;
        };
        if(body){
            myParameters['data'] = body;
        };
        const myResults = await axios(myParameters);
        
        return myResults.data;
    }

    async createUser(user, display, pass){
        return await this.requestConfig('/api/v1/users/', 'post', {
            body: {
            username: user,
            password: pass,
            displayName: display
        }})
    }

    async getUserInformation(userName){4
        const encodedUser = encodeURIComponent(userName)
        const response = await  this.requestConfig(`/api/v1/users/${encodedUser}`, 'get');
        return response;
    }1

    async userLogIn(username, password){
        const response = await this.requestConfig('/api/v1/users/login', 'post',{
            body: {
            username: username,
            password: password,
            
        }});
        this.userAuth = response['authToken'];
        console.log(this.userAuth);
        return response;
    }

    returnHeaders(){
        return ({'authorization':'Bearer ' + this.userAuth})
    }

    async getHighestBidForAll(){
        return await this.requestConfig('/api/v1/auctions/bids', 'get',{headers: this.returnHeaders});
     }

    async getHighestBidForOneItem(name){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/bids/${encodedItem}`, 'get', {
            headers: this.returnHeaders()
        });
    }

    async makeBid(name, bid){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/bids/${encodedItem}`, 'post', {
            body: {bidAmount: bid},
            headers: this.returnHeaders()
        });
    }
    
    async getAllItems(){
        console.log('entered axios')
        return await this.requestConfig(`/api/v1/auctions/items`, 'get',{
            headers: this.returnHeaders()})
            


            
    }

    async createNewItem(description, image, name){
        console.log(name);
        return await this.requestConfig('/api/v1/auctions/items', 'post', {
            body:
            {   description,
                image,
                name
            },
            headers: this.returnHeaders(),
        });
    }

    async deleteItem(name){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'delete',{
            headers: this.returnHeaders()
        });
    }

    async getSpecificItem(name){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'get', {
            headers: this.returnHeaders()
        });
    }

    async updateItem(description, image, name){
        const encodedItem = encodeURIComponent(name)
        return await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'put', {
            body: {
                description: description,
                image: image,
                name: name
            },
            headers:  this.returnHeaders()
        });
    }

    async establishWebSocket(){
        return await this.requestConfig('/api/v1/ws', get, {
            headers: this.returnHeaders()
        });
    }

}
export default new BiddrClient("http://localhost:8080")

window.axios = axios