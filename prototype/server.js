// ExpressJS Setup
const express = require('express');
const app = express();

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Constants
const PORT = 8080;
const HOST = '0.0.0.0';

// use static file
app.use(express.static(path.join(__dirname, 'views')));

// configure app to use body-parser
app.use(express.json());
app.use(express.urlencoded({ extended: false }));

// main page routing
app.get('/', (req, res)=>{
    res.sendFile(__dirname + '/index.html');
})

async function cc_call(fn_name, args){
    
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('medicalchannel');
    const contract = network.getContract('MedicalRecord');

    var result;
    
    if(fn_name == 'createRecordCopy') //증명서 정보 등록
    {   patno=args[0];
        patname=args[1];
        recordhash=args[2];
       result = await contract.submitTransaction('createRecordCopy', patno, patname, recordhash); 
    }
    else if(fn_name == 'queryPatNo') //환자 ID로 검색
        result = await contract.evaluateTransaction('queryPatNo', PatNo);
    else if(fn_name == 'queryTicket') // 티켓정보로 검색
        result = await contract.evaluateTransaction('queryTicket', TicketNumber);
    else if(fn_name == 'queryTX') //티켓정보로 트랜잭션 정보 검색
        result = await contract.evaluateTransaction('queryTX', TicketNumber);
    else
    
        result = 'not supported function'

    return result;
}

// create Record
app.post('/newRecord', async(req, res)=>{
    const No  = req.body.PatNo;
    const Name = req.body.PatName;
    const RecordHash = req.body.RecordHash;
    console.log("add Record PatNo: " + No);
    console.log("typeof PatNo: " + typeof(No));
    console.log("add Record PatName: " + Name);
    console.log("typeof PatName: " + typeof(Name));
    console.log("add Record Hash: " + RecordHash);
    console.log("typeof PatName: " + typeof(RecordHash));

    var args=[No, Name, RecordHash];
    result = cc_call('createRecordCopy', args)

    const myobj = {result: "success"}
    res.status(200).json(myobj) 
})

// find PatNo
app.post('/findPatNo', async (req,res)=>{
    const PatNo = req.body.PatNo;
    console.log("PatNo: " + req.body.PatNo);
    console.log("PatNo_typeof: " + typeof(PatNo));
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('medicalchannel');
    const contract = network.getContract('MedicalRecord');
    const result = await contract.evaluateTransaction('queryPatNo', PatNo);

    console.log("result :" + result);
    console.log(typeof(result)); 
    console.log(typeof(JSON.stringify(result))); 

    const myobj = JSON.parse(result)
    console.log("result :" + result);

    res.status(200).json(myobj)

});

// find PatTicket
app.post('/findTicketNumber', async (req,res)=>{
    const TicketNumber = req.body.TicketNumber;
    console.log("TicketNumber: " + req.body.TicketNumber);
    console.log("TicketNumber_typeof: " + typeof(TicketNumber));
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('medicalchannel');
    const contract = network.getContract('MedicalRecord');
    const result = await contract.evaluateTransaction('queryTicket', TicketNumber);
    console.log(JSON.stringify(result));
    const myobj = JSON.parse(result)
    res.status(200).json(myobj)
    // res.status(200).json(result)

});

app.post('/findtx', async (req,res)=>{
    const TicketNumber = req.body.TicketNumber;
    //console.log("TicketNumber: " + req.body.TicketNumber);
    //console.log("TicketNumber_typeof: " + typeof(TicketNumber));
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('medicalchannel');
    const contract = network.getContract('MedicalRecord');
    const result = await contract.evaluateTransaction('queryTX', TicketNumber);
    console.log("result : " + result)
    const myobj = JSON.parse(result)
    console.log("myobj : " + myobj)
    res.status(200).json(myobj)
    //res.status(200).json(result)
});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
