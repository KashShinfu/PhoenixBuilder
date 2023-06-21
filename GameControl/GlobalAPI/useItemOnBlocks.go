package GlobalAPI

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/mirror/chunk"
)

/*
让客户端点击 pos 处名为 blockName 且方块状态为 blockStates 的方块。
你可以对容器使用这样的操作，这会使得容器被打开
*/
func (g *GlobalAPI) ClickBlock(
	hotbarSlotID uint8,
	pos [3]int32,
	blockName string,
	blockStates map[string]interface{},
) error {
	standardRuntimeID, found := chunk.StateToRuntimeID(blockName, blockStates)
	if !found {
		return fmt.Errorf("ClickBlock: Failed to get the runtimeID of block %v; blockStates = %#v", blockName, blockStates)
	}
	blockRuntimeID := chunk.StandardRuntimeIDToNEMCRuntimeID(standardRuntimeID)
	if blockRuntimeID == chunk.AirRID || blockRuntimeID == chunk.NEMCAirRID {
		return fmt.Errorf("ClickBlock: Failed to converse StandardRuntimeID to NEMCRuntimeID; standardRuntimeID = %#v, blockName = %#v, blockStates = %#v", standardRuntimeID, blockName, blockStates)
	}
	// get block RunTime ID
	err := g.ChangeSelectedHotbarSlot(hotbarSlotID)
	if err != nil {
		return fmt.Errorf("ClickBlock: %v", err)
	}
	// change selected hotbar slot
	datas, err := g.Resources.Inventory.GetItemStackInfo(0, hotbarSlotID)
	if err != nil {
		return fmt.Errorf("ClickBlock: %v", err)
	}
	// get datas of the target item stack
	err = g.WritePacket(&packet.InventoryTransaction{
		LegacyRequestID:    0,
		LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
		Actions:            []protocol.InventoryAction{},
		TransactionData: &protocol.UseItemTransactionData{
			LegacyRequestID:    0,
			LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
			Actions:            []protocol.InventoryAction(nil),
			ActionType:         protocol.UseItemActionClickBlock,
			BlockPosition:      pos,
			HotBarSlot:         int32(hotbarSlotID),
			HeldItem:           datas,
			BlockRuntimeID:     blockRuntimeID,
		},
	})
	if err != nil {
		return fmt.Errorf("ClickBlock: %v", err)
	}
	err = g.WritePacket(&packet.PlayerAction{
		EntityRuntimeID: g.BotInfo.BotRunTimeID,
		ActionType:      protocol.PlayerActionStartBuildingBlock,
		BlockPosition:   pos,
	})
	if err != nil {
		return fmt.Errorf("ClickBlock: %v", err)
	}
	// send packet
	return nil
	// return
}