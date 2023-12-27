export namespace main {
	
	export class Message {
	    hash: string;
	    content: string;
	    name: string;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.content = source["content"];
	        this.name = source["name"];
	        this.timestamp = source["timestamp"];
	    }
	}

}

