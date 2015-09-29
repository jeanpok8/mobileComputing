package dht

import (
	"fmt"
    "testing"
	"strings"
	"math/big"
 )


type Contact struct{
	ip string
	port string
}

type DHTNode struct{
	id, nodeId  string 
	successor, predecessor *DHTNode
	contact  Contact
	finger   []*Fingers //links to Fingers structure 
	
}


/*this is the fingers structure, each DHTNode has finger which is populated
 *by fingers (i.e a start string and a pointer to a DHTNode, so the structure of the DHTNode)
 *will look like that:
 *
 *---> id:00 ip:nill port:nill
 *---> successor:01 predecessor:09 
 *--->finger[start,node],[start,node],[start,node]
 */
type Fingers struct {
	start string
	node *DHTNode
}


func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	if nodeId == nil { // if nodeId=nil, then generate generate it by calling generateNodeId fx.
		genNodeId := generateNodeId() 
		dhtNode.id = genNodeId
		dhtNode.contact.ip = ip 
	    dhtNode.contact.port = port
	    dhtNode.successor=dhtNode
	    dhtNode.predecessor=dhtNode
		dhtNode.finger=make([]*Fingers,160) // u can use either 3 or 160 bits(value of m).
		
		} else {
			
		dhtNode.id = *nodeId 
		dhtNode.contact.ip = ip 
	    dhtNode.contact.port = port
		dhtNode.successor=dhtNode
		dhtNode.predecessor=dhtNode
		dhtNode.finger=make([]*Fingers,160)// create a slice of type *Fingers with 160 fields.
	   
	}
	
	return dhtNode
}

/*Prnode=previous node, cnode=node that has to join the RING*/

 func (dhtNode *DHTNode)updateRing(cnode *DHTNode){
	fmt.Println("node id: ", cnode.id)
	if dhtNode.finger[0]==nil{
	
		for i:=1; i<=len(dhtNode.finger); i++{
			fingerID, _ :=calcFinger([]byte(dhtNode.id), i, len(dhtNode.finger))
			  if len(fingerID)<len(dhtNode.id){
				  fingerID=strings.Repeat("0",len(dhtNode.id)-len(fingerID))+fingerID
				  }
				  nodeT :=dhtNode.lookup(fingerID)
				  if nodeT.id!= fingerID{
					 nodeT=nodeT.successor
					  }
				  dhtNode.finger[i-1]= &Fingers{fingerID, nodeT}
				  fmt.Println(dhtNode.finger[i-1].node.id)
			  }
			  }
			  
			  for i :=1; i<=len(dhtNode.finger);i++{
				  fingerID,_:=calcFinger([]byte(cnode.id),i,len(dhtNode.finger))
				  if len(fingerID)<len(dhtNode.id){
				  fingerID=strings.Repeat("0",len(dhtNode.id)-len(fingerID))+fingerID
				  }
				  nodeT :=dhtNode.lookup(fingerID)
				  if nodeT.id!= fingerID{
					 nodeT=nodeT.successor
					  }
			   cnode.finger[i-1]=&Fingers{fingerID,nodeT}
	           fmt.Println(cnode.finger[i-1].node.id)
}

      node:=dhtNode.lookup(cnode.id)
      prnode:=node.successor 
      node.successor=cnode
      cnode.successor=prnode
      cnode.predecessor=node
	  prnode.predecessor=cnode
	  prnode.update()

}
 func (dhtNode *DHTNode) printRing() {
	 neNode:=dhtNode.successor
	 fmt.Println("id: ", dhtNode.id, "fingers: ", dhtNode.finger)
	 for neNode !=dhtNode{
		 fmt.Printf("id: %s fingers: ", neNode.id)
		 for i := 0; i < len(neNode.finger); i++ {
		 fmt.Printf("%s ", neNode.finger[i].node.id)
	 }
    fmt.Println()
	neNode=neNode.successor
 }
 }
 
 func (d *DHTNode) tostring() (out string) {
	out = "DHTNode{id: " + d.id + ", ip: " + d.contact.ip + ", port: " + d.contact.port + "}"
    
	return
}

func (d *DHTNode) lookup(hash string) *DHTNode {

	if between([]byte(d.id), []byte(d.successor.id), []byte(hash)) {
		return d
	}

	dist := distance(d.id, hash, len(d.finger))
	index := dist.BitLen() - 1
	if index < 0 {
		return d
	}
	fmt.Println("INDEX", index)

	// scroll down until your finger is not pointing at himself
	for ; index > 0 && d.finger[index].node == d; index-- {

	}
	// Viewing so we do not end up too far
	diff := big.Int{}
	diff.Sub(dist, distance(d.id, d.finger[index].node.id, len(d.finger)))
	for index > 0 && diff.Sign() < 0 {
		index--
		diff.Sub(dist, distance(d.id, d.finger[index].node.id, len(d.finger)))
	}
	//check so we do not point at ourselves
	if d.finger[index].node == d || diff.Sign() < 0 {
		fmt.Println("Error Error we are point our ourselves")
		return d.successor.lookup(hash)

	}

	return d.finger[index].node.lookup(hash)

	//	return d.successor.lookup(hash)
}

/*if s is in (one of) n fingers, update n's fingers with s*/

 func (dhtNode *DHTNode) update_table(s *DHTNode, i int){
    fmt.Println("updating finger", i, "on", dhtNode.id)
    if s.successor==dhtNode.finger[i-1].node{
	   dhtNode.finger[i-1].node=s
	   p :=dhtNode.predecessor
	   if p!=dhtNode{
		   p.update_table(s,i)
	  }
	}
 }
 func (dhtNode *DHTNode) update() {
	
	for i:=1; i<=len(dhtNode.finger);i++ {
	big_n:=big.Int{}
	sub_big_int:=big.Int{}	
	result := big.Int{}
	big_n.SetString(dhtNode.id, 16)
    sub_big_int.Exp(big.NewInt(2), big.NewInt(int64(i-1)), nil)
	result.Sub(&big_n, &sub_big_int)
	
	if result.Sign() < 0 {
			fmt.Println("Fix negative numbers")
			//will be used for 2^(nodes to be used)
			big_totalnodes := big.Int{}
			//the amount of nodes to be used
			//big_nodes := big.Int{}
			//used to do the calculation for sub
			big_negative := result

			//sets the nodes variable to a big int from the size of n.fingers

			big_totalnodes.Exp(big.NewInt(2), big.NewInt(int64(len(dhtNode.finger))), nil)
			//

			fmt.Println("the total number of nodes equals to: ", big_totalnodes)
			//calculate result
			fmt.Println("big_negative: ", big_negative)
			result.Add(&big_totalnodes, &big_negative)
            fmt.Println("This will be the final number !:")
			///// HERE, THERE MUST BE SENT SO THAT WE DO NOT TAKE -2 when it should be the node 7 for example
		}
	    bigString := fmt.Sprintf("%x", result.Bytes())
		fmt.Println(bigString)
		fmt.Println()
		fmt.Println()
		p := dhtNode.lookup(bigString)
		if p!= dhtNode {
		   p.update_table(dhtNode, i)
		    }
		  }
		}
		
	func (dhtNode *DHTNode) testCalcFingers(k int,m int){
	   bigN := big.Int{}
	   bigN.SetString(dhtNode.id, 16)
       fmt.Println(calcFinger(bigN.Bytes(), k, m))
	   }

func TestFinger160bits(t *testing.T) {
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	node1.updateRing(node2)
	node1.updateRing(node3)
	node1.updateRing(node4)
	node4.updateRing(node5)
	node3.updateRing(node6)
	node3.updateRing(node7)
	node3.updateRing(node8)
	node7.updateRing(node9)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	node3.testCalcFingers(0, 160)
	fmt.Println("")
	node3.testCalcFingers(1, 160)
	fmt.Println("")
	node3.testCalcFingers(80, 160)
	fmt.Println("")
	node3.testCalcFingers(120, 90)
	fmt.Println("")
	node3.testCalcFingers(160, 160)
	fmt.Println("")
}

/*func TestCreateRing(t *testing.T) { 
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"

	node0b := makeDHTNode(&id0, "localhost", "1111")
	node1b := makeDHTNode(&id1, "localhost", "1112")
	node2b := makeDHTNode(&id2, "localhost", "1113")
	node3b := makeDHTNode(&id3, "localhost", "1114")
	node4b := makeDHTNode(&id4, "localhost", "1115")
	node5b := makeDHTNode(&id5, "localhost", "1116")
	node6b := makeDHTNode(&id6, "localhost", "1117")
	node7b := makeDHTNode(&id7, "localhost", "1118")

	node0b.addToRing(node1b)
	node1b.addToRing(node2b)
	node1b.addToRing(node3b)
	node1b.addToRing(node4b)
	node4b.addToRing(node5b)
	node3b.addToRing(node6b)
	node3b.addToRing(node7b)

	fmt.Println("***RING STRUCTURE****")
	node1b.printRing()
	node3b.testCalcFingers(0, 3)
	node3b.testCalcFingers(1, 3)
	node3b.testCalcFingers(2, 3)
	node3b.testCalcFingers(3, 3)
}

func TestDHT2(t *testing.T) {
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	key1 := "2b230fe12d1c9c60a8e489d028417ac89de57635"
	key2 := "87adb987ebbd55db2c5309fd4b23203450ab0083"
	key3 := "74475501523a71c34f945ae4e87d571c2c57f6f3"

	fmt.Println("TEST: " + node1.lookup(key1).nodeId + " is responsible for " + key1)
	fmt.Println("TEST: " + node1.lookup(key2).nodeId + " is responsible for " + key2)
	fmt.Println("TEST: " + node1.lookup(key3).nodeId + " is responsible for " + key3)

	node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)
	node3.addToRing(node8)
	node7.addToRing(node9)

	fmt.Println("***RING STRUCTURE****")
	node1.printRing()

	nodeForKey1 := node1.lookup(key1)
	fmt.Println("dht node " + nodeForKey1.nodeId + " running at " + nodeForKey1.contact.ip + ":" + nodeForKey1.contact.port + " is responsible for " + key1)

	nodeForKey2 := node1.lookup(key2)
	fmt.Println("dht node " + nodeForKey2.nodeId + " running at " + nodeForKey2.contact.ip + ":" + nodeForKey2.contact.port + " is responsible for " + key2)

	nodeForKey3 := node1.lookup(key3)
	fmt.Println("dht node " + nodeForKey3.nodeId + " running at " + nodeForKey3.contact.ip + ":" + nodeForKey3.contact.port + " is responsible for " + key3)
}*/
