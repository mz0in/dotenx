import { Theme } from '@emotion/react'
import { atom, useAtom } from 'jotai'
import _ from 'lodash'
import { DragEventHandler, useEffect, useRef, useState } from 'react'
import {
	addEdge,
	ArrowHeadType,
	Connection,
	Edge,
	Elements,
	FlowElement,
	Node,
	OnLoadFunc,
	OnLoadParams,
	removeElements,
} from 'react-flow-renderer'
import { useQuery } from 'react-query'
import { getPipelineTriggers, PipelineData, QueryKey, TriggerData } from '../api'
import { InputOrSelectValue } from '../components/input-or-select'
import { EdgeData } from '../components/pipe-edge'
import { NodeType, TaskNodeData } from '../components/task-node'
import { Trigger } from '../containers/edge-settings'
import { selectedPipelineAtom } from '../containers/pipeline-select'
import { selectedPipelineDataAtom } from '../pages/home'
import { getLayoutedElements, NODE_HEIGHT, NODE_WIDTH } from './use-layout'

let id = 1
const getId = () => `task ${id++}`

let triggerId = 1
const getTriggerId = () => `trigger ${triggerId++}`

export const initialElements: Elements<TaskNodeData | EdgeData> = [
	{
		id: getId(),
		type: 'default',
		data: { name: 'task 1', type: '' },
		position: { x: 0, y: 0 },
	},
]

export const flowAtom = atom(initialElements)

export function useFlow() {
	const reactFlowWrapper = useRef<HTMLDivElement>(null)
	const [reactFlowInstance, setReactFlowInstance] = useState<OnLoadParams | null>(null)
	const [elements, setElements] = useAtom(flowAtom)
	const [pipeline] = useAtom(selectedPipelineAtom)
	const [selectedPipelineData] = useAtom(selectedPipelineDataAtom)

	const triggersQuery = useQuery(
		[QueryKey.GetPipelineTriggers, selectedPipelineData?.name],
		() => {
			if (selectedPipelineData) return getPipelineTriggers(selectedPipelineData.name)
		},
		{ enabled: !!selectedPipelineData }
	)

	useEffect(() => {
		if (!pipeline) return
		const elements = mapPipelineToElements(pipeline)
		const triggers = mapTriggersToElements(triggersQuery.data?.data)
		const layout = getLayoutedElements(
			[...elements, ...triggers],
			'TB',
			NODE_WIDTH,
			NODE_HEIGHT
		)
		setElements(layout)
	}, [pipeline, setElements, triggersQuery.data])

	const onConnect = (params: Edge | Connection) => {
		setElements((els) => addEdge({ ...params, arrowHeadType: ArrowHeadType.Arrow }, els))
	}

	const onElementsRemove = (elementsToRemove: Elements) => {
		setElements((els) => removeElements(elementsToRemove, els))
	}

	const onLoad: OnLoadFunc = (reactFlowInstance) => {
		setReactFlowInstance(reactFlowInstance)
		reactFlowInstance.fitView()
	}

	const onDragOver: DragEventHandler<HTMLDivElement> = (event) => {
		event.preventDefault()
		event.dataTransfer.dropEffect = 'move'
	}

	const onDrop: DragEventHandler<HTMLDivElement> = (event) => {
		event.preventDefault()

		const reactFlowBounds = reactFlowWrapper.current?.getBoundingClientRect()
		const type = event.dataTransfer.getData('application/reactflow')
		if (!reactFlowInstance || !reactFlowBounds) return
		const position = reactFlowInstance.project({
			x: event.clientX - reactFlowBounds.left,
			y: event.clientY - reactFlowBounds.top,
		})
		const id = type === NodeType.Default ? getId() : getTriggerId()
		const newNode: FlowElement<TaskNodeData> = {
			id,
			type,
			position,
			data: { name: id, type: '' },
		}

		setElements((es) => es.concat(newNode))
	}

	const updateElement = (id: string, data: TaskNodeData | EdgeData) => {
		setElements((els) => els.map((el) => (el.id === id ? { ...el, data } : el)))
	}

	return {
		reactFlowWrapper,
		elements,
		onConnect,
		onElementsRemove,
		onLoad,
		onDragOver,
		onDrop,
		updateElement,
	}
}

export function getNodeColor(theme: Theme, node: Node) {
	switch (node.type) {
		case NodeType.Default:
			return theme.color.text
		case NodeType.Trigger:
			return theme.color.primary
		default:
			return theme.color.text
	}
}

function mapPipelineToElements(pipeline: PipelineData): Elements<TaskNodeData | EdgeData> {
	const nodes = Object.entries(pipeline.manifest.tasks).map(([key, value]) => {
		const bodyEntries = _.toPairs(value.body).map(([fieldName, fieldValue]) => {
			let inputOrSelectValue = { type: 'text', data: '' } as InputOrSelectValue
			if (typeof fieldValue === 'string') {
				inputOrSelectValue = { type: 'text', data: fieldValue }
			} else {
				inputOrSelectValue = {
					type: 'option',
					data: fieldValue.key,
					groupName: fieldValue.source,
					iconUrl: '',
				}
			}
			return [fieldName, inputOrSelectValue]
		})
		const body = _.fromPairs(bodyEntries)
		console.log(body)
		return {
			id: key,
			position: { x: 0, y: 0 },
			type: NodeType.Default,
			data: {
				name: key,
				type: value.type,
				integration: value.integration,
				...body,
			},
		}
	})
	const edges = Object.entries(pipeline.manifest.tasks).flatMap(([target, task]) =>
		Object.entries(task.executeAfter).map(([source, triggers]) => ({
			id: `${source}to${target}`,
			source,
			target,
			arrowHeadType: ArrowHeadType.Arrow,
			data: {
				triggers: triggers as Trigger[],
			},
		}))
	)

	return [...nodes, ...edges]
}

function mapTriggersToElements(triggers: TriggerData[] | undefined) {
	if (!triggers) return []

	const triggerNodes = triggers.map((trigger) => ({
		id: trigger.name,
		position: { x: 0, y: 0 },
		type: NodeType.Trigger,
		data: trigger,
	}))

	return triggerNodes
}
