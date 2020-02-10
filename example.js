const fs = require('fs')

const cmd = process.argv[2]

let inputString = '';

process.stdin.on('data', inputStdin => {
    inputString += inputStdin;
});

process.stdin.on('end', async _ => {
    await main();    
});

async function main() {
	fs.writeFileSync('debug.txt', inputString)
	const resource = JSON.parse(inputString);
	fs.writeFileSync('debug.txt', JSON.stringify({
		id: resource && resource.id,
		data: process.argv,
		raw: inputString,
		stdin: resource,
	}, null, '\t'))

	let ret;
	switch(cmd) {
	case 'read':
		ret = await read(resource)
		break
	case 'create':
		ret = await create(resource)
		break
	case 'update':
		ret = await update(resource)
		break
	case 'delete':
		ret = await delete_resource(resource)
		break
	default:
		throw new Error(`Unknown cmd ${cmd}`)
	}

	console.log(JSON.stringify(ret, null, '\t'))
	process.exit(0)
}

async function read(resource) {
	let json = require(`./tmp/${resource.ID}.json`)
	return {
		ID: json.ID,
		...json,
		DeepObject: JSON.stringify(json)
	}
}

async function create(resource) {
	let obj = {
		ID: (+new Date).toString(16),
		...JSON.parse(resource.Payload),
		aaa: "bbb",
		debug: {
			abc: 123
		}
	}
	fs.writeFileSync(`./tmp/${obj.ID}.json`, JSON.stringify(obj, null, '\t'))
	return read(obj)
}

async function update(resource) {
	let obj = {
		ID: resource.ID,
		...JSON.parse(resource.Payload),
		aaa: "bbb",
		debug: {
			abc: 123
		}
	}
	fs.writeFileSync(`./tmp/${resource.ID}.json`, JSON.stringify(obj, null, '\t'))
	return read(resource)
}

async function delete_resource(resource) {
	fs.unlinkSync(`./tmp/${resource.ID}.json`)
	return {}
}