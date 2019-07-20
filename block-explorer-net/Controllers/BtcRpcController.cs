using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Threading.Tasks;
using block_explorer_net.Models;
using Microsoft.AspNetCore.Mvc;
using NBitcoin;
using NBitcoin.RPC;

namespace block_explorer_net.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class BtcRpcController : ControllerBase
    {
        private readonly RPCClient _client;
        private readonly Network _network;
        
        public BtcRpcController()
        {
            _network = Network.RegTest;

            _client = new RPCClient(
                credentials: new NetworkCredential(userName: "myuser", password: "SomeDecentp4ssw0rd"), 
                address: new Uri("http://localhost:18443"),
                network: _network);
        }
        
        [HttpGet("block")]
        public async Task<List<BlockDto>> GetBlockHashes([FromQuery] int count)
        {
            var blockCount = await _client.GetBlockCountAsync();

            List<BlockDto> hashList = new List<BlockDto>(blockCount);
            for (int i = 0; i < count && i < blockCount; i++)
            {
                var hash = await _client.GetBlockAsync(blockCount-i);
                hashList.Add(new BlockDto(hash));
            }

            return hashList;
        }

        [HttpGet("block/{hash}")]
        public async Task<BlockDto> GetBlock([FromRoute] string hash)
        {
            var block = await _client.GetBlockAsync(uint256.Parse(hash));
            return new BlockDto(block);
        }

        [HttpGet("tx/{txid}")]
        public async Task<TransactionDto> GetTx([FromRoute] string txid)
        {
            var tx = await _client.GetRawTransactionAsync(uint256.Parse(txid));
            return new TransactionDto(tx);
        }
    
        [HttpPost] 
        public async Task Send(string address, long satoshi) 
        {
            var btcAddress = BitcoinAddress.Create(address, _network);
            var amount = new Money(satoshis: satoshi);

            await _client.SendToAddressAsync(btcAddress, amount);
        }
    }
}
