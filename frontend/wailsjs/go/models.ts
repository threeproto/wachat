export namespace params {
	
	export class Message {
	    id: number;
	    hash: string;
	    content: string;
	    name: string;
	    timestamp: number;
	    is_stored: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.hash = source["hash"];
	        this.content = source["content"];
	        this.name = source["name"];
	        this.timestamp = source["timestamp"];
	        this.is_stored = source["is_stored"];
	    }
	}

}

