using System;
using System.Collections.Generic;
using System.Linq;
using NBitcoin;

namespace block_explorer_net.Models
{
    public class BlockDto
    {
        public BlockDto(){}
        public BlockDto(Block block){
            this.Hash = block.GetHash().ToString();
            this.HashPrevBlock = block.Header.HashPrevBlock.ToString();
            this.BlockTime = block.Header.BlockTime;
            this.Transactions = block.Transactions.Select( x=> new TransactionDto(x)).ToList();
        }

        public string Hash { get; }
        public string HashPrevBlock { get; }
        public DateTimeOffset BlockTime { get; }
        public List<TransactionDto> Transactions { get; set; }
    }
}