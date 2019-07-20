using System;
using NBitcoin;

namespace block_explorer_net.Models
{
    public class TransactionDto
    {
        public TransactionDto(){}
        public TransactionDto(Transaction tx)
        {
            this.Hash = tx.GetHash().ToString();
        }

        public string Hash { get; set;}
    }
}