# HyperLedger 仮想通貨・お買い物依頼・共同購入 chaincode（仮）

## 概要

[HyperLedger/fabric](https://github.com/hyperledger/fabric)で公開されているサンプルを元に、「仮想通貨のやり取り」「お買い物依頼」「共同購入」を実現するchaincodeです。       
（時間が無いため、細かなエラー処理は実装できてません…。）       


## 各chaincode概要

### [cc1（仮想通貨）](https://github.com/moritaku-92/go-chaincode/tree/master/skillgroup_cc1)

仮想通貨のやり取りを行うことが出来るchaincodeです。      
keyに「口座名義人」、valueに「口座残高」を設定しています。      
振込や残高確認、ユーザの追加・削除など基本的な機能を備えています。


### [cc2（お買い物依頼）](https://github.com/moritaku-92/go-chaincode/tree/master/skillgroup_cc2)

お買い物依頼が出来るchaincodeです。      
keyに「依頼番号」、valueに「依頼内容（json）」を設定しています。      
依頼を達成時に依頼者から受注者に報酬が支払われます。報酬のやり取りはcc1を使用します。      
お買い物依頼、依頼受領、依頼キャンセル、依頼一覧取得などの機能を備えています。     


### [cc3（共同購入）](https://github.com/moritaku-92/go-chaincode/tree/master/skillgroup_cc3)

共同購入（クラウドファンディングのイメージ）が出来るchaincodeです。  
keyに「共同購入番号」、valueに「共同購入の内容（json）」を設定しています。    
共同購入者数が一定数に達した際に、支援者から指定の金額が減算されます。金額のやり取りはcc1を使用します。       
共同購入の提案、共同購入への参加などの機能を備えています。


